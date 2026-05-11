#!/usr/bin/env python3
"""
he_cli.py — HE language CLI
Usage:
    python he_cli.py run genesis.he --out MyGame
"""

import os, sys, re, subprocess, platform, shutil
from collections import deque

# -----------------------
# Config / type map
# -----------------------
TYPE_MAP = {"int": "int", "dec": "float", "string": "std::string"}
FETCHED = set()

# -----------------------
# Tokenizer
# -----------------------
TOKEN_SPEC = [
    ("NUMBER",   r"-?\d+(\.\d+)?"),
    ("STRING",   r'"([^"\\]|\\.)*"'),
    ("IDENT",    r"[A-Za-z_][A-Za-z0-9_]*"),
    ("OP",       r"==|!=|<=|>=|[:=+\-]"),
    ("PUNC",     r"[{}\[\]()]",),
    ("SKIP",     r"[ \t]+"),
    ("NEWLINE",  r"\n"),
    ("COMMENT",  r"//[^\n]*"),
    ("MISMATCH", r"."),
]

TOK_REGEX = re.compile("|".join(f"(?P<{n}>{p})" for n, p in TOKEN_SPEC))

class Token:
    def __init__(self, typ, val, lineno, col):
        self.type = typ
        self.val = val
        self.lineno = lineno
        self.col = col
    def __repr__(self):
        return f"Token({self.type!s},{self.val!r},{self.lineno}:{self.col})"

def tokenize(text, filename="<stdin>"):
    """Return list of Token objects. Removes BOM if present."""
    if text.startswith("\ufeff"):
        text = text[1:]
    tokens = []
    lineno = 1
    line_start = 0
    for mo in TOK_REGEX.finditer(text):
        kind = mo.lastgroup
        val = mo.group(0)
        col = mo.start() - line_start + 1
        if kind == "NEWLINE":
            lineno += 1
            line_start = mo.end()
            continue
        if kind in ("SKIP", "COMMENT"):
            continue
        if kind == "STRING":
            # remove surrounding quotes, keep contents as-is (no escape processing needed now)
            val = val[1:-1]
        tokens.append(Token(kind, val, lineno, col))
    tokens.append(Token("EOF", "", lineno, 1))
    return tokens

# -----------------------
# AST
# -----------------------
class Program:
    def __init__(self):
        self.objects = []

class ObjectNode:
    def __init__(self, name):
        self.name = name
        self.sections = []  # list of PropertiesSection

class PropertiesSection:
    def __init__(self, name):
        self.name = name
        self.props = []  # list of Property

class Property:
    def __init__(self, name, type_, value):
        self.name = name
        self.type_ = type_
        self.value = value

# -----------------------
# Parser
# -----------------------
class ParserError(Exception):
    pass

class Parser:
    def __init__(self, tokens, filename="<stdin>"):
        self.toks = tokens
        self.i = 0
        self.filename = filename

    def cur(self):
        return self.toks[self.i]

    def eat(self, typ=None, val=None):
        tok = self.cur()
        if typ and tok.type != typ:
            raise ParserError(f"{self.filename}:{tok.lineno}:{tok.col}: Expected token type {typ}, got {tok.type} ('{tok.val}')")
        if val and tok.val != val:
            raise ParserError(f"{self.filename}:{tok.lineno}:{tok.col}: Expected token '{val}', got '{tok.val}'")
        self.i += 1
        return tok

    def accept(self, typ, val=None):
        tok = self.cur()
        if tok.type == typ and (val is None or tok.val == val):
            self.i += 1
            return tok
        return None

    # Top-level program
    def parse_program(self):
        prog = Program()
        while self.cur().type != "EOF":
            # fetch "file.he"
            if self.accept("IDENT", "fetch"):
                # fetch expects a STRING token (filename in quotes)
                s = self.eat("STRING")
                # actual processing of fetch is handled outside parser (process_file)
                continue

            # make <Name> { ... }
            if self.accept("IDENT", "make"):
                obj = self.parse_object()
                prog.objects.append(obj)
                continue

            tok = self.cur()
            raise ParserError(f"{self.filename}:{tok.lineno}:{tok.col}: Unexpected token {tok.val}")
        return prog

    # parse_object: assumes 'make' already consumed
    # grammar:
    # make Identifier { <ObjectBody> }
    # ObjectBody ::= Identifier "has" [ ":" ] { PropertySection | nested make ... }
    def parse_object(self):
        name_tok = self.eat("IDENT")
        name = name_tok.val
        self.eat("PUNC", "{")
        obj = ObjectNode(name)

        # inside object body: expect lines like:
        # <ObjectName> has:   (or without colon)
        #     Section [
        #         prop: type = value
        #     ]
        while not self.accept("PUNC", "}"):
            # Expect object name repetition for section header
            objname_tok = self.eat("IDENT")
            objname = objname_tok.val
            if objname != name:
                raise ParserError(f"{self.filename}:{objname_tok.lineno}:{objname_tok.col}: Expected '{name} has:' (found '{objname}')")

            # next should be 'has'
            self.eat("IDENT", "has")
            # optional colon
            self.accept("OP", ":")

            # now zero or more property sections (each starts with IDENT then '[')
            while self.cur().type == "IDENT":
                sec = self.parse_properties_section()
                obj.sections.append(sec)

        return obj

    # parse_properties_section: SectionName [ prop: type = value ... ]
    # property values can be STRING token or multiple IDENT/NUMBER tokens on the same source line
    def parse_properties_section(self):
        sec_name_tok = self.eat("IDENT")
        sec_name = sec_name_tok.val
        # expect '['
        self.eat("PUNC", "[")
        sec = PropertiesSection(sec_name)

        while not self.accept("PUNC", "]"):
            # property: Identifier ':' Type '=' Value
            prop_name_tok = self.eat("IDENT")
            prop_name = prop_name_tok.val
            self.eat("OP", ":")
            type_tok = self.eat("IDENT")
            type_name = type_tok.val

            # placeholder: support "Name: type-" (uninitialized) if next token is OP '-'
            if self.accept("OP", "-"):
                value = None
                sec.props.append(Property(prop_name, type_name, value))
                continue

            self.eat("OP", "=")

            # Collect value tokens that belong to the same source line (so "Hello World" across two IDENT tokens works)
            parts = []
            start_line = self.cur().lineno
            while True:
                tok = self.cur()
                if tok.type in ("IDENT", "NUMBER", "STRING"):
                    parts.append(self.eat(tok.type).val)
                    # stop if next token is on a different source line or is PUNC ']' (end of section) or next property start
                    nxt = self.cur()
                    if nxt.type == "PUNC" and nxt.val == "]":
                        break
                    if nxt.lineno != start_line:
                        break
                    # allow multiple identifiers/numbers on same line
                    continue
                else:
                    # If a non-value token appears (e.g. PUNC or EOF), stop
                    break

            value = " ".join(parts) if parts else None
            sec.props.append(Property(prop_name, type_name, value))

        return sec

# -----------------------
# File processing (with fetch)
# -----------------------
def read_file(path):
    with open(path, "r", encoding="utf-8") as f:
        return f.read()

def process_file(path):
    """
    Parse a given file and return a Program AST.
    Handles recursive fetch by scanning for fetch "..." and merging fetched programs.
    """
    if path in FETCHED:
        return Program()
    if not os.path.isfile(path):
        raise FileNotFoundError(path)
    FETCHED.add(path)

    text = read_file(path)
    tokens = tokenize(text, path)
    parser = Parser(tokens, path)
    prog = parser.parse_program()

    # handle fetch directives (scan the source; easier than making parser build import AST nodes)
    for m in re.finditer(r'\bfetch\s+"([^"]+)"', text):
        fname = m.group(1)
        # support relative paths
        base_dir = os.path.dirname(path) or "."
        fetched_path = os.path.join(base_dir, fname)
        if os.path.isfile(fetched_path):
            nested_prog = process_file(fetched_path)
            # merge objects from nested program into current prog
            prog.objects.extend(nested_prog.objects)
        else:
            raise FileNotFoundError(f"fetched module not found: {fetched_path}")

    return prog

def merge_programs(programs):
    final = Program()
    for p in programs:
        final.objects.extend(p.objects)
    return final

# -----------------------
# C++ generation
# -----------------------
def cpp_literal_for(type_name, val):
    """Create a C++ literal for the given value and type."""
    if val is None:
        return '0' if type_name != "std::string" else '""'
    # if value already looks numeric
    if re.match(r"^-?\d+(\.\d+)?$", str(val)):
        return str(val)
    # else it's a string-like; escape double quotes and backslashes
    s = str(val).replace("\\", "\\\\").replace('"', '\\"')
    return f'"{s}"'

def generate_cpp(prog):
    lines = []
    lines.append("#include <string>")
    lines.append("#include <iostream>")
    lines.append("")
    lines.append("using namespace std;")
    lines.append("")

    for obj in prog.objects:
        lines.append(f"struct {obj.name} " + "{")
        lines.append("public:")
        for sec in obj.sections:
            for prop in sec.props:
                if prop.type_ == "end":
                    # end flag -> bool <name>_alive
                    lines.append(f"    bool {prop.name}_alive = true;")
                else:
                    ctype = TYPE_MAP.get(prop.type_, "int")
                    lit = cpp_literal_for("std::string" if ctype == "std::string" else ctype, prop.value)
                    lines.append(f"    {ctype} {prop.name} = {lit};")
        lines.append("};")
        lines.append("")

    # main: instantiate first object and no-op
    if prog.objects:
        main_obj = prog.objects[0]
        lines.append("int main() {")
        lines.append(f"    {main_obj.name} obj;")
        lines.append("    return 0;")
        lines.append("}")

    return "\n".join(lines)

# -----------------------
# Compilation helpers
# -----------------------
def find_compiler():
    if shutil.which("g++"):
        return "g++"
    if shutil.which("clang++"):
        return "clang++"
    if shutil.which("cl"):
        return "cl"
    return None

def compile_cpp(cpp_path, exe_name):
    compiler = find_compiler()
    if compiler is None:
        raise RuntimeError("No C++ compiler found (g++/clang++/cl). Install one or add to PATH.")
    system = platform.system()
    if compiler in ("g++", "clang++"):
        # On Windows with MinGW the exe_name should end with .exe if desired; otherwise keep as provided.
        cmd = [compiler, "-O2", "-std=c++17", cpp_path, "-o", exe_name]
        print("Compiling:", " ".join(cmd))
        subprocess.check_call(cmd)
        return exe_name + (".exe" if system == "Windows" and not exe_name.lower().endswith(".exe") else "")
    elif compiler == "cl":
        # MSVC needs to be run in Developer Command Prompt with environment prepared
        exe_file = exe_name if exe_name.lower().endswith(".exe") else exe_name + ".exe"
        cmd = ["cl", "/EHsc", "/std:c++17", "/O2", cpp_path, "/Fe:" + exe_file]
        print("Compiling (MSVC):", " ".join(cmd))
        subprocess.check_call(" ".join(cmd), shell=True)
        return exe_file

# -----------------------
# CLI entry
# -----------------------
def run_he(entry, outname):
    prog = process_file(entry)
    cpp_code = generate_cpp(prog)
    cpp_path = "output.cpp"
    with open(cpp_path, "w", encoding="utf-8") as f:
        f.write(cpp_code)
    print("Transpiled →", cpp_path)
    exe = compile_cpp(cpp_path, outname)
    print("Executable →", exe)

if __name__ == "__main__":
    import argparse
    ap = argparse.ArgumentParser()
    ap.add_argument("run")
    ap.add_argument("entry")
    ap.add_argument("--out", "-o", default="game")
    args = ap.parse_args()

    run_he(args.entry, args.out)
