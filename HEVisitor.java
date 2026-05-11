// Generated from HE.g4 by ANTLR 4.13.0
import org.antlr.v4.runtime.tree.ParseTreeVisitor;

/**
 * This interface defines a complete generic visitor for a parse tree produced
 * by {@link HEParser}.
 *
 * @param <T> The return type of the visit operation. Use {@link Void} for
 * operations with no return type.
 */
public interface HEVisitor<T> extends ParseTreeVisitor<T> {
	/**
	 * Visit a parse tree produced by {@link HEParser#program}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitProgram(HEParser.ProgramContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#line}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitLine(HEParser.LineContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#summon}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitSummon(HEParser.SummonContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#withAssets}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitWithAssets(HEParser.WithAssetsContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#assetList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAssetList(HEParser.AssetListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#asset}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAsset(HEParser.AssetContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#assetType}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAssetType(HEParser.AssetTypeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#object}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitObject(HEParser.ObjectContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#objectBody}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitObjectBody(HEParser.ObjectBodyContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#section}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitSection(HEParser.SectionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#properties}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitProperties(HEParser.PropertiesContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#abilities}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAbilities(HEParser.AbilitiesContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#reactions}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitReactions(HEParser.ReactionsContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#memories}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitMemories(HEParser.MemoriesContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#trigger}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitTrigger(HEParser.TriggerContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#block}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitBlock(HEParser.BlockContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#property}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitProperty(HEParser.PropertyContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#action}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAction(HEParser.ActionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#paramList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitParamList(HEParser.ParamListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#parameter}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitParameter(HEParser.ParameterContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#statement}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitStatement(HEParser.StatementContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#say}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitSay(HEParser.SayContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#change}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitChange(HEParser.ChangeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#decide}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitDecide(HEParser.DecideContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#repeat}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitRepeat(HEParser.RepeatContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#call}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitCall(HEParser.CallContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#waitStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitWaitStmt(HEParser.WaitStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#returnStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitReturnStmt(HEParser.ReturnStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#argList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitArgList(HEParser.ArgListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#expression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitExpression(HEParser.ExpressionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#logicOr}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitLogicOr(HEParser.LogicOrContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#logicAnd}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitLogicAnd(HEParser.LogicAndContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#comparison}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitComparison(HEParser.ComparisonContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#arithmetic}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitArithmetic(HEParser.ArithmeticContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#term}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitTerm(HEParser.TermContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#factor}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitFactor(HEParser.FactorContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#unary}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitUnary(HEParser.UnaryContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#primary}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitPrimary(HEParser.PrimaryContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#callExpression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitCallExpression(HEParser.CallExpressionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#arrayLiteral}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitArrayLiteral(HEParser.ArrayLiteralContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#type}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitType(HEParser.TypeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#baseType}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitBaseType(HEParser.BaseTypeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#globalStatement}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitGlobalStatement(HEParser.GlobalStatementContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#comment}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitComment(HEParser.CommentContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#compOp}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitCompOp(HEParser.CompOpContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#addOp}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAddOp(HEParser.AddOpContext ctx);
	/**
	 * Visit a parse tree produced by {@link HEParser#multOp}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitMultOp(HEParser.MultOpContext ctx);
}