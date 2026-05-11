// Generated from src\main\antlr\HE_New.g4 by ANTLR 4.13.0
package he_new.parser;

package he_new.parser;

import org.antlr.v4.runtime.tree.ParseTreeVisitor;

/**
 * This interface defines a complete generic visitor for a parse tree produced
 * by {@link HE_NewParser}.
 *
 * @param <T> The return type of the visit operation. Use {@link Void} for
 * operations with no return type.
 */
public interface HE_NewVisitor<T> extends ParseTreeVisitor<T> {
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#program}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitProgram(HE_NewParser.ProgramContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#line}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitLine(HE_NewParser.LineContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#summon}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitSummon(HE_NewParser.SummonContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#withAssets}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitWithAssets(HE_NewParser.WithAssetsContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#assetList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAssetList(HE_NewParser.AssetListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#asset}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAsset(HE_NewParser.AssetContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#assetType}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAssetType(HE_NewParser.AssetTypeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#object}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitObject(HE_NewParser.ObjectContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#objectBody}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitObjectBody(HE_NewParser.ObjectBodyContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#objectSection}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitObjectSection(HE_NewParser.ObjectSectionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#propertiesSection}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitPropertiesSection(HE_NewParser.PropertiesSectionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#propertyBlock}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitPropertyBlock(HE_NewParser.PropertyBlockContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#property}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitProperty(HE_NewParser.PropertyContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#abilitiesSection}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAbilitiesSection(HE_NewParser.AbilitiesSectionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#abilityBlock}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAbilityBlock(HE_NewParser.AbilityBlockContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#method}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitMethod(HE_NewParser.MethodContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#paramList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitParamList(HE_NewParser.ParamListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#parameter}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitParameter(HE_NewParser.ParameterContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#reactionSection}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitReactionSection(HE_NewParser.ReactionSectionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#type}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitType(HE_NewParser.TypeContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#statement}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitStatement(HE_NewParser.StatementContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#printStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitPrintStmt(HE_NewParser.PrintStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#callStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitCallStmt(HE_NewParser.CallStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#showStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitShowStmt(HE_NewParser.ShowStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#waitStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitWaitStmt(HE_NewParser.WaitStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#assignmentStmt}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAssignmentStmt(HE_NewParser.AssignmentStmtContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#argList}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitArgList(HE_NewParser.ArgListContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#expression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitExpression(HE_NewParser.ExpressionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#additiveExpression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitAdditiveExpression(HE_NewParser.AdditiveExpressionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#multiplicativeExpression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitMultiplicativeExpression(HE_NewParser.MultiplicativeExpressionContext ctx);
	/**
	 * Visit a parse tree produced by {@link HE_NewParser#primaryExpression}.
	 * @param ctx the parse tree
	 * @return the visitor result
	 */
	T visitPrimaryExpression(HE_NewParser.PrimaryExpressionContext ctx);
}