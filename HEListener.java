// Generated from HE.g4 by ANTLR 4.13.0
import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link HEParser}.
 */
public interface HEListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link HEParser#program}.
	 * @param ctx the parse tree
	 */
	void enterProgram(HEParser.ProgramContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#program}.
	 * @param ctx the parse tree
	 */
	void exitProgram(HEParser.ProgramContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#line}.
	 * @param ctx the parse tree
	 */
	void enterLine(HEParser.LineContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#line}.
	 * @param ctx the parse tree
	 */
	void exitLine(HEParser.LineContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#summon}.
	 * @param ctx the parse tree
	 */
	void enterSummon(HEParser.SummonContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#summon}.
	 * @param ctx the parse tree
	 */
	void exitSummon(HEParser.SummonContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#withAssets}.
	 * @param ctx the parse tree
	 */
	void enterWithAssets(HEParser.WithAssetsContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#withAssets}.
	 * @param ctx the parse tree
	 */
	void exitWithAssets(HEParser.WithAssetsContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#assetList}.
	 * @param ctx the parse tree
	 */
	void enterAssetList(HEParser.AssetListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#assetList}.
	 * @param ctx the parse tree
	 */
	void exitAssetList(HEParser.AssetListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#asset}.
	 * @param ctx the parse tree
	 */
	void enterAsset(HEParser.AssetContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#asset}.
	 * @param ctx the parse tree
	 */
	void exitAsset(HEParser.AssetContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#assetType}.
	 * @param ctx the parse tree
	 */
	void enterAssetType(HEParser.AssetTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#assetType}.
	 * @param ctx the parse tree
	 */
	void exitAssetType(HEParser.AssetTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#object}.
	 * @param ctx the parse tree
	 */
	void enterObject(HEParser.ObjectContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#object}.
	 * @param ctx the parse tree
	 */
	void exitObject(HEParser.ObjectContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#objectBody}.
	 * @param ctx the parse tree
	 */
	void enterObjectBody(HEParser.ObjectBodyContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#objectBody}.
	 * @param ctx the parse tree
	 */
	void exitObjectBody(HEParser.ObjectBodyContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#section}.
	 * @param ctx the parse tree
	 */
	void enterSection(HEParser.SectionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#section}.
	 * @param ctx the parse tree
	 */
	void exitSection(HEParser.SectionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#properties}.
	 * @param ctx the parse tree
	 */
	void enterProperties(HEParser.PropertiesContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#properties}.
	 * @param ctx the parse tree
	 */
	void exitProperties(HEParser.PropertiesContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#abilities}.
	 * @param ctx the parse tree
	 */
	void enterAbilities(HEParser.AbilitiesContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#abilities}.
	 * @param ctx the parse tree
	 */
	void exitAbilities(HEParser.AbilitiesContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#reactions}.
	 * @param ctx the parse tree
	 */
	void enterReactions(HEParser.ReactionsContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#reactions}.
	 * @param ctx the parse tree
	 */
	void exitReactions(HEParser.ReactionsContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#memories}.
	 * @param ctx the parse tree
	 */
	void enterMemories(HEParser.MemoriesContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#memories}.
	 * @param ctx the parse tree
	 */
	void exitMemories(HEParser.MemoriesContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#trigger}.
	 * @param ctx the parse tree
	 */
	void enterTrigger(HEParser.TriggerContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#trigger}.
	 * @param ctx the parse tree
	 */
	void exitTrigger(HEParser.TriggerContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#block}.
	 * @param ctx the parse tree
	 */
	void enterBlock(HEParser.BlockContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#block}.
	 * @param ctx the parse tree
	 */
	void exitBlock(HEParser.BlockContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#property}.
	 * @param ctx the parse tree
	 */
	void enterProperty(HEParser.PropertyContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#property}.
	 * @param ctx the parse tree
	 */
	void exitProperty(HEParser.PropertyContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#action}.
	 * @param ctx the parse tree
	 */
	void enterAction(HEParser.ActionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#action}.
	 * @param ctx the parse tree
	 */
	void exitAction(HEParser.ActionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#paramList}.
	 * @param ctx the parse tree
	 */
	void enterParamList(HEParser.ParamListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#paramList}.
	 * @param ctx the parse tree
	 */
	void exitParamList(HEParser.ParamListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#parameter}.
	 * @param ctx the parse tree
	 */
	void enterParameter(HEParser.ParameterContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#parameter}.
	 * @param ctx the parse tree
	 */
	void exitParameter(HEParser.ParameterContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#statement}.
	 * @param ctx the parse tree
	 */
	void enterStatement(HEParser.StatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#statement}.
	 * @param ctx the parse tree
	 */
	void exitStatement(HEParser.StatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#say}.
	 * @param ctx the parse tree
	 */
	void enterSay(HEParser.SayContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#say}.
	 * @param ctx the parse tree
	 */
	void exitSay(HEParser.SayContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#change}.
	 * @param ctx the parse tree
	 */
	void enterChange(HEParser.ChangeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#change}.
	 * @param ctx the parse tree
	 */
	void exitChange(HEParser.ChangeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#decide}.
	 * @param ctx the parse tree
	 */
	void enterDecide(HEParser.DecideContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#decide}.
	 * @param ctx the parse tree
	 */
	void exitDecide(HEParser.DecideContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#repeat}.
	 * @param ctx the parse tree
	 */
	void enterRepeat(HEParser.RepeatContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#repeat}.
	 * @param ctx the parse tree
	 */
	void exitRepeat(HEParser.RepeatContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#call}.
	 * @param ctx the parse tree
	 */
	void enterCall(HEParser.CallContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#call}.
	 * @param ctx the parse tree
	 */
	void exitCall(HEParser.CallContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#waitStmt}.
	 * @param ctx the parse tree
	 */
	void enterWaitStmt(HEParser.WaitStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#waitStmt}.
	 * @param ctx the parse tree
	 */
	void exitWaitStmt(HEParser.WaitStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#returnStmt}.
	 * @param ctx the parse tree
	 */
	void enterReturnStmt(HEParser.ReturnStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#returnStmt}.
	 * @param ctx the parse tree
	 */
	void exitReturnStmt(HEParser.ReturnStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#argList}.
	 * @param ctx the parse tree
	 */
	void enterArgList(HEParser.ArgListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#argList}.
	 * @param ctx the parse tree
	 */
	void exitArgList(HEParser.ArgListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterExpression(HEParser.ExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitExpression(HEParser.ExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#logicOr}.
	 * @param ctx the parse tree
	 */
	void enterLogicOr(HEParser.LogicOrContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#logicOr}.
	 * @param ctx the parse tree
	 */
	void exitLogicOr(HEParser.LogicOrContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#logicAnd}.
	 * @param ctx the parse tree
	 */
	void enterLogicAnd(HEParser.LogicAndContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#logicAnd}.
	 * @param ctx the parse tree
	 */
	void exitLogicAnd(HEParser.LogicAndContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#comparison}.
	 * @param ctx the parse tree
	 */
	void enterComparison(HEParser.ComparisonContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#comparison}.
	 * @param ctx the parse tree
	 */
	void exitComparison(HEParser.ComparisonContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#arithmetic}.
	 * @param ctx the parse tree
	 */
	void enterArithmetic(HEParser.ArithmeticContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#arithmetic}.
	 * @param ctx the parse tree
	 */
	void exitArithmetic(HEParser.ArithmeticContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#term}.
	 * @param ctx the parse tree
	 */
	void enterTerm(HEParser.TermContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#term}.
	 * @param ctx the parse tree
	 */
	void exitTerm(HEParser.TermContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#factor}.
	 * @param ctx the parse tree
	 */
	void enterFactor(HEParser.FactorContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#factor}.
	 * @param ctx the parse tree
	 */
	void exitFactor(HEParser.FactorContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#unary}.
	 * @param ctx the parse tree
	 */
	void enterUnary(HEParser.UnaryContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#unary}.
	 * @param ctx the parse tree
	 */
	void exitUnary(HEParser.UnaryContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#primary}.
	 * @param ctx the parse tree
	 */
	void enterPrimary(HEParser.PrimaryContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#primary}.
	 * @param ctx the parse tree
	 */
	void exitPrimary(HEParser.PrimaryContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#callExpression}.
	 * @param ctx the parse tree
	 */
	void enterCallExpression(HEParser.CallExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#callExpression}.
	 * @param ctx the parse tree
	 */
	void exitCallExpression(HEParser.CallExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#arrayLiteral}.
	 * @param ctx the parse tree
	 */
	void enterArrayLiteral(HEParser.ArrayLiteralContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#arrayLiteral}.
	 * @param ctx the parse tree
	 */
	void exitArrayLiteral(HEParser.ArrayLiteralContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#type}.
	 * @param ctx the parse tree
	 */
	void enterType(HEParser.TypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#type}.
	 * @param ctx the parse tree
	 */
	void exitType(HEParser.TypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#baseType}.
	 * @param ctx the parse tree
	 */
	void enterBaseType(HEParser.BaseTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#baseType}.
	 * @param ctx the parse tree
	 */
	void exitBaseType(HEParser.BaseTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#globalStatement}.
	 * @param ctx the parse tree
	 */
	void enterGlobalStatement(HEParser.GlobalStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#globalStatement}.
	 * @param ctx the parse tree
	 */
	void exitGlobalStatement(HEParser.GlobalStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#comment}.
	 * @param ctx the parse tree
	 */
	void enterComment(HEParser.CommentContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#comment}.
	 * @param ctx the parse tree
	 */
	void exitComment(HEParser.CommentContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#compOp}.
	 * @param ctx the parse tree
	 */
	void enterCompOp(HEParser.CompOpContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#compOp}.
	 * @param ctx the parse tree
	 */
	void exitCompOp(HEParser.CompOpContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#addOp}.
	 * @param ctx the parse tree
	 */
	void enterAddOp(HEParser.AddOpContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#addOp}.
	 * @param ctx the parse tree
	 */
	void exitAddOp(HEParser.AddOpContext ctx);
	/**
	 * Enter a parse tree produced by {@link HEParser#multOp}.
	 * @param ctx the parse tree
	 */
	void enterMultOp(HEParser.MultOpContext ctx);
	/**
	 * Exit a parse tree produced by {@link HEParser#multOp}.
	 * @param ctx the parse tree
	 */
	void exitMultOp(HEParser.MultOpContext ctx);
}