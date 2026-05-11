import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.tree.*;

public class Test {
    public static void main(String[] args) throws Exception {
        System.out.println("Testing HE Parser...");
        
        // Test with proper line endings
        String input = "summon \"core\" as system\n" +
                       "print \"Hello\"\n";
        
        // Create lexer and parser
        CharStream stream = CharStreams.fromString(input);
        HELexer lexer = new HELexer(stream);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        HEParser parser = new HEParser(tokens);
        
        // Parse
        HEParser.ProgramContext tree = parser.program();
        
        System.out.println("✅ Parse successful!");
        System.out.println("Tree: " + tree.toStringTree(parser));
        
        // Print tokens for debugging
        System.out.println("\nTokens:");
        tokens.fill();
        for (Token token : tokens.getTokens()) {
            if (token.getType() != Token.EOF) {
                String name = HELexer.VOCABULARY.getSymbolicName(token.getType());
                System.out.println("  " + name + ": '" + token.getText() + "'");
            }
        }
    }
}