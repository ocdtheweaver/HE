// src/main/kotlin/Main.kt
import ast.Program
import compiler.KotlinCompiler
import visitor.AstVisitor
import org.antlr.v4.runtime.CharStreams
import org.antlr.v4.runtime.CommonTokenStream
import he_new.parser.HE_NewLexer
import he_new.parser.HE_NewParser
import java.io.File
// Add this to your Main.kt or a separate CommentExtractor class
import org.antlr.v4.runtime.*
import org.antlr.v4.runtime.tree.*

fun extractComments(tokens: CommonTokenStream): List<Pair<Int, String>> {
    val comments = mutableListOf<Pair<Int, String>>()
    
    for (token in tokens.tokens) {
        if (token.channel == Token.HIDDEN_CHANNEL && token.type == HE_NewLexer.COMMENT) {
            val text = token.text.removeSurrounding("~").trim()
            comments.add(token.tokenIndex to text)
        }
    }
    
    return comments
}

fun main(args: Array<String>) {
    if (args.isEmpty()) {
        println("Usage: HE_NewCompiler <input-file.he>")
        return
    }
    
    val inputFile = File(args[0])
    if (!inputFile.exists()) {
        println("File not found: ${inputFile.path}")
        return
    }
    
    try {
        // Read and parse input
        val input = CharStreams.fromFileName(inputFile.path)
        val lexer = HE_NewLexer(input)
        val tokens = CommonTokenStream(lexer)
        val parser = HE_NewParser(tokens)
        
        val tree = parser.program()
        
        // Convert to AST
        val visitor = AstVisitor()
        val program = visitor.visitProgram(tree) as Program
        
        // Generate Kotlin code
        val compiler = KotlinCompiler()
        val kotlinCode = compiler.compile(program)
        
        // Write output
        val outputFile = File(inputFile.nameWithoutExtension + ".kt")
        outputFile.writeText(kotlinCode)
        
        println("Successfully compiled ${inputFile.name} to ${outputFile.name}")
        
    } catch (e: Exception) {
        println("Error: ${e.message}")
        e.printStackTrace()
    }
}