// src/main/kotlin/ast/AstNodes.kt
package ast

sealed class Node

sealed class Statement : Node()
data class PrintStmt(val expression: Expression) : Statement()
data class CallStmt(val target: String, val method: String, val args: List<Expression>) : Statement()
data class ShowStmt(val asset: String, val centered: Boolean) : Statement()
data class WaitStmt(val duration: Expression) : Statement()
data class AssignmentStmt(val name: String, val value: Expression) : Statement()

sealed class Expression : Node()
data class StringLiteral(val value: String) : Expression()
data class Identifier(val name: String) : Expression()
data class NumberLiteral(val value: Double) : Expression()
data class BinaryExpression(val left: Expression, val operator: String, val right: Expression) : Expression()

data class Program(val lines: List<Line>) : Node()

sealed class Line : Node()
data class SummonLine(val type: String, val name: String) : Line()
data class WithAssetsLine(val assets: List<Asset>) : Line()
data class ObjectDef(val name: String, val parent: String?, val body: ObjectBody) : Line()
// CommentLine removed since comments are in HIDDEN channel

data class Asset(val type: AssetType, val path: String, val name: String?)
data class ObjectBody(val sections: List<ObjectSection>)

sealed class ObjectSection : Node()
data class PropertiesSection(val name: String, val properties: List<Property>) : ObjectSection()
data class AbilitiesSection(val name: String, val methods: List<Method>) : ObjectSection()
data class ReactionSection(val event: String, val statements: List<Statement>) : ObjectSection()

data class Property(val name: String, val value: Expression)
data class Method(val name: String, val params: List<Parameter>, val body: List<Statement>)
data class Parameter(val name: String, val type: Type)

enum class AssetType { IMAGE, SOUND, MUSIC, VIDEO, FONT, SHADER }
enum class Type { STRING, NUMBER, BOOLEAN }