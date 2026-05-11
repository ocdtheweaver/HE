import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.tree.*;
import java.io.*;

public class RealTest {
    public static void main(String[] args) throws IOException {
        System.out.println("Testing real HE example...");
        
        CharStream input = CharStreams.fromFileName("real_test.he");
        HELexer lexer = new HELexer(input);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        HEParser parser = new HEParser(tokens);
        
        // Add error listener to see detailed errors
        parser.removeErrorListeners();
        parser.addErrorListener(new BaseErrorListener() {
            @Override
            public void syntaxError(Recognizer<?, ?> recognizer, Object offendingSymbol,
                    int line, int charPositionInLine, String msg, RecognitionException e) {
                System.out.println("Syntax error at line " + line + ":" + charPositionInLine + " " + msg);
            }
        });
        
        HEParser.ProgramContext tree = parser.program();
        System.out.println("✅ Real example parse successful!");
        System.out.println("Tree depth: " + tree.toStringTree(parser).length() + " chars");
    }
}