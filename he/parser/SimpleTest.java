package he.parser;

import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.tree.*;

public class SimpleTest {
    public static void main(String[] args) throws Exception {
        System.out.println("Testing HE Parser...");
        
        String input = "summon \"core\" as system\nprint \"Hello\"";
        
        CharStream stream = CharStreams.fromString(input);
        HELexer lexer = new HELexer(stream);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        HEParser parser = new HEParser(tokens);
        
        HEParser.ProgramContext tree = parser.program();
        System.out.println("✅ Parse successful!");
        System.out.println("Tree: " + tree.toStringTree(parser));
    }
}