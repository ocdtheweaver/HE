// Generated from src\main\antlr\HE_New.g4 by ANTLR 4.13.0
package he_new.parser;

package he_new.parser;

import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link HE_NewParser}.
 */
public interface HE_NewListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#program}.
	 * @param ctx the parse tree
	 */
	void enterProgram(HE_NewParser.ProgramContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#program}.
	 * @param ctx the parse tree
	 */
	void exitProgram(HE_NewParser.ProgramContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#line}.
	 * @param ctx the parse tree
	 */
	void enterLine(HE_NewParser.LineContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#line}.
	 * @param ctx the parse tree
	 */
	void exitLine(HE_NewParser.LineContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#summon}.
	 * @param ctx the parse tree
	 */
	void enterSummon(HE_NewParser.SummonContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#summon}.
	 * @param ctx the parse tree
	 */
	void exitSummon(HE_NewParser.SummonContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#withAssets}.
	 * @param ctx the parse tree
	 */
	void enterWithAssets(HE_NewParser.WithAssetsContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#withAssets}.
	 * @param ctx the parse tree
	 */
	void exitWithAssets(HE_NewParser.WithAssetsContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#assetList}.
	 * @param ctx the parse tree
	 */
	void enterAssetList(HE_NewParser.AssetListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#assetList}.
	 * @param ctx the parse tree
	 */
	void exitAssetList(HE_NewParser.AssetListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#asset}.
	 * @param ctx the parse tree
	 */
	void enterAsset(HE_NewParser.AssetContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#asset}.
	 * @param ctx the parse tree
	 */
	void exitAsset(HE_NewParser.AssetContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#assetType}.
	 * @param ctx the parse tree
	 */
	void enterAssetType(HE_NewParser.AssetTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#assetType}.
	 * @param ctx the parse tree
	 */
	void exitAssetType(HE_NewParser.AssetTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#object}.
	 * @param ctx the parse tree
	 */
	void enterObject(HE_NewParser.ObjectContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#object}.
	 * @param ctx the parse tree
	 */
	void exitObject(HE_NewParser.ObjectContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#objectBody}.
	 * @param ctx the parse tree
	 */
	void enterObjectBody(HE_NewParser.ObjectBodyContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#objectBody}.
	 * @param ctx the parse tree
	 */
	void exitObjectBody(HE_NewParser.ObjectBodyContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#objectSection}.
	 * @param ctx the parse tree
	 */
	void enterObjectSection(HE_NewParser.ObjectSectionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#objectSection}.
	 * @param ctx the parse tree
	 */
	void exitObjectSection(HE_NewParser.ObjectSectionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#propertiesSection}.
	 * @param ctx the parse tree
	 */
	void enterPropertiesSection(HE_NewParser.PropertiesSectionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#propertiesSection}.
	 * @param ctx the parse tree
	 */
	void exitPropertiesSection(HE_NewParser.PropertiesSectionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#propertyBlock}.
	 * @param ctx the parse tree
	 */
	void enterPropertyBlock(HE_NewParser.PropertyBlockContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#propertyBlock}.
	 * @param ctx the parse tree
	 */
	void exitPropertyBlock(HE_NewParser.PropertyBlockContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#property}.
	 * @param ctx the parse tree
	 */
	void enterProperty(HE_NewParser.PropertyContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#property}.
	 * @param ctx the parse tree
	 */
	void exitProperty(HE_NewParser.PropertyContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#abilitiesSection}.
	 * @param ctx the parse tree
	 */
	void enterAbilitiesSection(HE_NewParser.AbilitiesSectionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#abilitiesSection}.
	 * @param ctx the parse tree
	 */
	void exitAbilitiesSection(HE_NewParser.AbilitiesSectionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#abilityBlock}.
	 * @param ctx the parse tree
	 */
	void enterAbilityBlock(HE_NewParser.AbilityBlockContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#abilityBlock}.
	 * @param ctx the parse tree
	 */
	void exitAbilityBlock(HE_NewParser.AbilityBlockContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#method}.
	 * @param ctx the parse tree
	 */
	void enterMethod(HE_NewParser.MethodContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#method}.
	 * @param ctx the parse tree
	 */
	void exitMethod(HE_NewParser.MethodContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#paramList}.
	 * @param ctx the parse tree
	 */
	void enterParamList(HE_NewParser.ParamListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#paramList}.
	 * @param ctx the parse tree
	 */
	void exitParamList(HE_NewParser.ParamListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#parameter}.
	 * @param ctx the parse tree
	 */
	void enterParameter(HE_NewParser.ParameterContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#parameter}.
	 * @param ctx the parse tree
	 */
	void exitParameter(HE_NewParser.ParameterContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#reactionSection}.
	 * @param ctx the parse tree
	 */
	void enterReactionSection(HE_NewParser.ReactionSectionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#reactionSection}.
	 * @param ctx the parse tree
	 */
	void exitReactionSection(HE_NewParser.ReactionSectionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#type}.
	 * @param ctx the parse tree
	 */
	void enterType(HE_NewParser.TypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#type}.
	 * @param ctx the parse tree
	 */
	void exitType(HE_NewParser.TypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#statement}.
	 * @param ctx the parse tree
	 */
	void enterStatement(HE_NewParser.StatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#statement}.
	 * @param ctx the parse tree
	 */
	void exitStatement(HE_NewParser.StatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#printStmt}.
	 * @param ctx the parse tree
	 */
	void enterPrintStmt(HE_NewParser.PrintStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#printStmt}.
	 * @param ctx the parse tree
	 */
	void exitPrintStmt(HE_NewParser.PrintStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#callStmt}.
	 * @param ctx the parse tree
	 */
	void enterCallStmt(HE_NewParser.CallStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#callStmt}.
	 * @param ctx the parse tree
	 */
	void exitCallStmt(HE_NewParser.CallStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#showStmt}.
	 * @param ctx the parse tree
	 */
	void enterShowStmt(HE_NewParser.ShowStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#showStmt}.
	 * @param ctx the parse tree
	 */
	void exitShowStmt(HE_NewParser.ShowStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#waitStmt}.
	 * @param ctx the parse tree
	 */
	void enterWaitStmt(HE_NewParser.WaitStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#waitStmt}.
	 * @param ctx the parse tree
	 */
	void exitWaitStmt(HE_NewParser.WaitStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#assignmentStmt}.
	 * @param ctx the parse tree
	 */
	void enterAssignmentStmt(HE_NewParser.AssignmentStmtContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#assignmentStmt}.
	 * @param ctx the parse tree
	 */
	void exitAssignmentStmt(HE_NewParser.AssignmentStmtContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#argList}.
	 * @param ctx the parse tree
	 */
	void enterArgList(HE_NewParser.ArgListContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#argList}.
	 * @param ctx the parse tree
	 */
	void exitArgList(HE_NewParser.ArgListContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterExpression(HE_NewParser.ExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitExpression(HE_NewParser.ExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#additiveExpression}.
	 * @param ctx the parse tree
	 */
	void enterAdditiveExpression(HE_NewParser.AdditiveExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#additiveExpression}.
	 * @param ctx the parse tree
	 */
	void exitAdditiveExpression(HE_NewParser.AdditiveExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#multiplicativeExpression}.
	 * @param ctx the parse tree
	 */
	void enterMultiplicativeExpression(HE_NewParser.MultiplicativeExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#multiplicativeExpression}.
	 * @param ctx the parse tree
	 */
	void exitMultiplicativeExpression(HE_NewParser.MultiplicativeExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link HE_NewParser#primaryExpression}.
	 * @param ctx the parse tree
	 */
	void enterPrimaryExpression(HE_NewParser.PrimaryExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link HE_NewParser#primaryExpression}.
	 * @param ctx the parse tree
	 */
	void exitPrimaryExpression(HE_NewParser.PrimaryExpressionContext ctx);
}