#!/usr/bin/env python3
# hec.py — HE transpiler (hybrid typing + JS backend + runtime checks)
# Usage:
#   python hec.py build genesis.he --target js --out genesis.js [--runtime-checks]
#   python hec.py run genesis.he [--runtime-checks]  (will run node if available)
#
# Generates JS that runs in Node or the browser.

import re, os, sys, argparse, subprocess
from shutil import which

# -----------------------
# AST classes
# -----------------------
class Property:
    def __init__(self, name, type_, value=None):
        self.name = name
        self.type_ = type_
        self.value = value

class Action:
    def __init__(self, name, params=None):
        self.name = name
        # params: list of tuples (name, type_or_None)
        self.params = params or []
        self.statements = []

class Block:
    def __init__(self, name):
        self.name = name
        self.properties = []
        self.actions = []
        self.base_path = ""

class ObjectNode:
    def __init__(self, name):
        self.name = name
        self.blocks = []

# -----------------------
# Utilities
# -----------------------
def error(msg, filename, lineno):
    raise SyntaxError(f"{filename}:{lineno}: {msg}")

def is_number_literal(s):
    return re.match(r'^-?\d+(\.\d+)?$', s) is not None

def quote_string_if_needed(s):
    s = s.strip()
    if s.startswith('"') and s.endswith('"'):
        return s
    # if looks like a number or identifier, keep raw
    if is_number_literal(s) or re.match(r'^[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*$', s):
        return s
    # otherwise quote
    return '"' + s.replace('"', '\\"') + '"'

# -----------------------
# Parser
# -----------------------
def parse_he(lines, filename="<stdin>"):
    current_obj = None
    current_block = None
    current_action = None
    stack = []
    root_objects = []
    line_num = 0

    while line_num < len(lines):
        raw = lines[line_num]
        line = raw.strip()
        line_num += 1
        if not line or line.startswith("//"):
            continue

        # fetch "file.he"  (imports)
        m = re.match(r'fetch\s+"([^"]+)"', line)
        if m:
            modpath = m.group(1)
            if not os.path.isfile(modpath):
                error(f"Module '{modpath}' not found", filename, line_num)
            with open(modpath, "r", encoding="utf-8") as f:
                imported = f.readlines()
            imported_ast = parse_he(imported, filename=modpath)
            root_objects.extend(imported_ast)
            continue

        # make Name {
        m = re.match(r"make\s+([A-Za-z_][A-Za-z0-9_]*)\s*{", line)
        if m:
            name = m.group(1)
            obj = ObjectNode(name)
            current_obj = obj
            stack.append(obj)
            root_objects.append(obj)
            continue

        # Close object: }
        if line == "}":
            if not stack:
                error("Unexpected '}'", filename, line_num)
            stack.pop()
            current_obj = None
            current_block = None
            current_action = None
            continue

        # Block header: Identifier has:
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s+has\s*:', line)
        if m and current_obj:
            block_name = m.group(1)
            block = Block(block_name)
            current_obj.blocks.append(block)
            current_block = block
            stack.append(block)
            continue

        # Block start: Name [
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*\[', line)
        if m and current_obj and not current_block:
            block = Block(m.group(1))
            current_obj.blocks.append(block)
            current_block = block
            stack.append(block)
            continue

        # Generic '[' open (we accept it)
        if line == "[":
            continue

        # Action start with optional typed params: name(param: type, p2)
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*(?:\(([^)]*)\))?\s*\[', line)
        if m and current_block:
            name = m.group(1)
            raw_params = m.group(2) or ""
            params = []
            for p in [x.strip() for x in raw_params.split(",") if x.strip()]:
                m2 = re.match(r'^([A-Za-z_][A-Za-z0-9_]*)(?:\s*:\s*(int|dec|string))?$', p)
                if not m2:
                    error(f"Invalid parameter syntax: '{p}'", filename, line_num)
                pname = m2.group(1)
                ptype = m2.group(2)
                params.append((pname, ptype))
            action = Action(name, params)
            current_block.actions.append(action)
            current_action = action
            stack.append(action)
            continue

        # with = "path"
        m = re.match(r'with\s*=\s*"([^"]+)"', line)
        if m and current_block and not current_action:
            current_block.base_path = m.group(1)
            continue

        # shorthand property: name with = "file";
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*(?:with\s*=\s*"([^"]+)")\s*;?', line)
        if m and current_block and not current_action:
            name = m.group(1)
            val = (current_block.base_path or "") + (m.group(2) or "")
            current_block.properties.append(Property(name, "string", val))
            continue

        # typed property: name: type = value
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*:\s*(int|dec|string)\s*=\s*(.+)', line)
        if m and current_block and not current_action:
            name, t, val = m.group(1), m.group(2), m.group(3).strip()
            if val.endswith(";"):
                val = val[:-1].strip()
            current_block.properties.append(Property(name, t, val))
            continue

        # placeholder property: name: type-
        m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*:\s*(int|dec|string)\s*-', line)
        if m and current_block and not current_action:
            current_block.properties.append(Property(m.group(1), m.group(2), None))
            continue

        # end flag
        m = re.match(r'end\s+([A-Za-z_][A-Za-z0-9_]*)', line)
        if m and current_block:
            current_block.properties.append(Property(m.group(1), "end", None))
            continue

        # close block or action: ]
        if line == "]":
            if not stack:
                error("Unexpected ']'", filename, line_num)
            popped = stack.pop()
            if isinstance(popped, Action):
                current_action = None
                # set current_block to nearest Block in stack
                current_block = None
                for s in reversed(stack):
                    if isinstance(s, Block):
                        current_block = s
                        break
            elif isinstance(popped, Block):
                current_block = None
                current_action = None
            else:
                current_block = None
                current_action = None
            continue

        # statements inside action
        if current_action:
            stmt = line
            # print
            m = re.match(r'print\s+(.+)', stmt)
            if m:
                current_action.statements.append(("print", m.group(1).strip()))
                continue
            # assignment: a = value  or a = Obj.prop
            m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+)', stmt)
            if m:
                left = m.group(1)
                right = m.group(2).strip()
                current_action.statements.append(("assign", left, right))
                continue
            # increase/reduce
            m = re.match(r'(increase|reduce)\s+([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)?)\s+by\s+(.+)', stmt)
            if m:
                op, target, by = m.group(1), m.group(2), m.group(3).strip()
                current_action.statements.append(("op", op, target, by))
                continue
            # if (expr)
            m = re.match(r'if\s*\((.+)\)', stmt)
            if m:
                current_action.statements.append(("if", m.group(1).strip()))
                continue
            m = re.match(r'otherwise\s*\((.+)\)', stmt)
            if m:
                current_action.statements.append(("elseif", m.group(1).strip()))
                continue
            if stmt == "orelse":
                current_action.statements.append(("else", None))
                continue
            # short loops: repeat name until cond OR run name until cond
            m = re.match(r'(repeat|run)\s+([A-Za-z_][A-Za-z0-9_]*)\s+until\s+(.+)', stmt)
            if m:
                current_action.statements.append(("loop_short", m.group(1), m.group(2), m.group(3).strip()))
                continue
            # call: name()
            m = re.match(r'([A-Za-z_][A-Za-z0-9_]*)\s*\(\s*\)', stmt)
            if m:
                current_action.statements.append(("call", m.group(1)))
                continue
            # bind ... with ... so ...
            m = re.match(r'bind\s+(.+?)\s+with\s+(.+?)\s+so\s+(.+)', stmt)
            if m:
                current_action.statements.append(("bind", m.group(1).strip(), m.group(2).strip(), m.group(3).strip()))
                continue
            # unrecognized -> raw
            current_action.statements.append(("raw", stmt))
            continue

        # else unexpected
        error("Invalid or unexpected line: " + line, filename, line_num)

    return root_objects

# -----------------------
# JS generation (uses type info)
# -----------------------
def js_type_check_snippet(varname, typ):
    if typ == "int":
        return f'if (typeof {varname} !== "number" || !Number.isInteger({varname})) throw new TypeError("Expected int for {varname}");'
    if typ == "dec":
        return f'if (typeof {varname} !== "number") throw new TypeError("Expected dec (number) for {varname}");'
    if typ == "string":
        return f'if (typeof {varname} !== "string") throw new TypeError("Expected string for {varname}");'
    return ""

def transform_expr_js(expr):
    # replace 'and'/'or'
    expr = expr.replace(" and ", " && ").replace(" or ", " || ")
    # replace identifier sequences with this.<id> where appropriate
    # but allow dotted access like Obj.prop (we convert bare identifiers to this.<id>)
    def repl(m):
        tok = m.group(0)
        # if token is a number or quoted string, keep
        if is_number_literal(tok) or (tok.startswith('"') and tok.endswith('"')):
            return tok
        # keep JS reserved boolean true/false
        if tok in ("true", "false", "null"):
            return tok
        # If token contains a dot it's probably Obj.prop — convert to this.Obj.prop? We'll leave as-is.
        return "this." + tok
    # only replace bare identifiers (words)
    return re.sub(r'[A-Za-z_][A-Za-z0-9_]*', repl, expr)

def cond_to_js(cond):
    cond = cond.strip()
    m = re.match(r'end\s+([A-Za-z_][A-Za-z0-9_]*)', cond)
    if m:
        return f'this.{m.group(1)}_alive'
    # Replace selector syntax A|B to this.A_B_alive if used; support simple comparators
    cond = cond.replace("|", ".")
    return transform_expr_js(cond)

def value_to_js(val):
    v = val.strip()
    if v.startswith('"') and v.endswith('"'):
        return v
    if is_number_literal(v):
        return v
    # identifier or dotted path
    if '.' in v:
        # if it's like Obj.prop -> treat as this.Obj.prop? user probably meant object.property
        parts = v.split('.')
        if len(parts) == 2:
            return f"this.{parts[0]}.{parts[1]}"
        return v
    return f"this.{v}"

def generate_js(objects, runtime_checks=False):
    out = []
    out.append("// Generated by hec.py — JavaScript backend")
    out.append("(function(){")
    out.append("  // Runtime checks: " + str(bool(runtime_checks)))
    out.append("")
    for obj in objects:
        out.append(f"  class {obj.name} " + "{")
        out.append("    constructor() {")
        # properties from blocks
        for block in obj.blocks:
            for prop in block.properties:
                if prop.type_ == "end":
                    out.append(f"      this.{prop.name}_alive = true;")
                else:
                    if prop.value is None:
                        default = "0" if prop.type_ in ("int","dec") else '""'
                    else:
                        val = prop.value.strip()
                        if prop.type_ == "string" and not (val.startswith('"') and val.endswith('"')):
                            val = '"' + val + '"'
                        default = val
                    out.append(f"      this.{prop.name} = {default};")
        out.append("    }")
        out.append("")
        # actions
        for block in obj.blocks:
            for action in block.actions:
                params = ", ".join([pname for (pname, ptype) in action.params])
                out.append(f"    {action.name}({params}) " + "{")
                # runtime checks for params
                if runtime_checks:
                    for (pname, ptype) in action.params:
                        if ptype:
                            snippet = js_type_check_snippet(pname, ptype)
                            if snippet:
                                out.append("      " + snippet)
                # body
                for stmt in action.statements:
                    kind = stmt[0]
                    if kind == "print":
                        expr = stmt[1]
                        # support printing identifiers -> this.<id>
                        out.append(f"      console.log({value_to_js(expr)});")
                        out.append(f"      document.body.innerHTML += '<pre>' + String({value_to_js(expr)}) + '</pre>';")
                    elif kind == "assign":
                        left, right = stmt[1], stmt[2]
                        out.append(f"      this.{left} = {value_to_js(right)};")
                    elif kind == "op":
                        op, target, by = stmt[1], stmt[2], stmt[3]
                        tgt_js = target if '.' in target else f"this.{target}"
                        op_js = "-=" if op == "reduce" else "+="
                        out.append(f"      {tgt_js} {op_js} {value_to_js(by)};")
                    elif kind == "if":
                        out.append("      if (" + transform_expr_js(stmt[1]) + ") {")
                    elif kind == "elseif":
                        out.append("      } else if (" + transform_expr_js(stmt[1]) + ") {")
                    elif kind == "else":
                        out.append("      } else {")
                    elif kind == "loop_short":
                        looptype, actionname, cond = stmt[1], stmt[2], stmt[3]
                        cond_js = cond_to_js(cond)
                        if looptype == "repeat":
                            out.append(f"      do {{ this.{actionname}(); }} while (!({cond_js}));")
                        else:
                            out.append(f"      while (!({cond_js})) {{ this.{actionname}(); }}")
                    elif kind == "call":
                        out.append(f"      this.{stmt[1]}();")
                    elif kind == "bind":
                        a,b,so = stmt[1], stmt[2], stmt[3]
                        out.append(f"      {a} = {b}; // then {so}")
                    else:
                        out.append(f"      // raw: {stmt[1] if len(stmt)>1 else stmt}")
                # close any hanging ifs — we won't auto-close; assume user balanced
                out.append("    }")
                out.append("")
        out.append("  }")
        out.append("")
    # auto-main: instantiate first object and call first action if present
    if objects:
        main = objects[0]
        # find first action
        first_action = None
        for b in main.blocks:
            if b.actions:
                first_action = b.actions[0].name
                break
        out.append(f"  const __main = new {main.name}();")
        if first_action:
            out.append(f"  if (typeof __main.{first_action} === 'function') __main.{first_action}();")
    out.append("})();")
    return "\n".join(out)

# -----------------------
# CLI & run
# -----------------------
def build_file(path, out="output.js", runtime_checks=False):
    with open(path, "r", encoding="utf-8") as f:
        lines = f.readlines()
    ast = parse_he(lines, filename=path)
    js = generate_js(ast, runtime_checks=runtime_checks)
    with open(out, "w", encoding="utf-8") as fo:
        fo.write(js)
    print(f"Wrote {out}")
    return out

def run_file(path, runtime_checks=False):
    out = build_file(path, out="output.js", runtime_checks=runtime_checks)
    if which("node"):
        subprocess.run(["node", out])
    else:
        print(f"Built {out}. Open it in a browser (see index.html) or install Node to run it.")

def main():
    ap = argparse.ArgumentParser(description="HE compiler (transpile to JS).")
    ap.add_argument("cmd", choices=["build", "run"], help="build or run")
    ap.add_argument("file", help=".he source")
    ap.add_argument("--out", default="output.js", help="output file")
    ap.add_argument("--runtime-checks", action="store_true", help="emit runtime type checks for typed params")
    args = ap.parse_args()
    if args.cmd == "build":
        build_file(args.file, out=args.out, runtime_checks=args.runtime_checks)
    else:
        run_file(args.file, runtime_checks=args.runtime_checks)

if __name__ == "__main__":
    main()
