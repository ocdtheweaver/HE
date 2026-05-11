// Generated from src\main\antlr\HE_New.g4 by ANTLR 4.13.0
package he_new.parser;

package he_new.parser;

import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class HE_NewParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.0", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		SUMMON=1, WITH=2, NAMED=3, AS=4, MAKE=5, LIKE=6, HAS=7, CAN=8, ON=9, IS=10, 
		PRINT=11, TELL=12, TO=13, SHOW=14, CENTERED=15, WAIT=16, SECONDS=17, AND=18, 
		IMAGE=19, SOUND=20, MUSIC=21, VIDEO=22, FONT=23, SHADER=24, STRING_TYPE=25, 
		NUMBER_TYPE=26, BOOLEAN_TYPE=27, COLON=28, LBRACK=29, RBRACK=30, LPAREN=31, 
		RPAREN=32, COMMA=33, TILDA=34, LBRACE=35, RBRACE=36, DOT=37, PLUS=38, 
		MINUS=39, MULT=40, DIV=41, COMMENT=42, ID=43, NUMBER=44, STRING=45, NEWLINE=46, 
		WS=47;
	public static final int
		RULE_program = 0, RULE_line = 1, RULE_summon = 2, RULE_withAssets = 3, 
		RULE_assetList = 4, RULE_asset = 5, RULE_assetType = 6, RULE_object = 7, 
		RULE_objectBody = 8, RULE_objectSection = 9, RULE_propertiesSection = 10, 
		RULE_propertyBlock = 11, RULE_property = 12, RULE_abilitiesSection = 13, 
		RULE_abilityBlock = 14, RULE_method = 15, RULE_paramList = 16, RULE_parameter = 17, 
		RULE_reactionSection = 18, RULE_type = 19, RULE_statement = 20, RULE_printStmt = 21, 
		RULE_callStmt = 22, RULE_showStmt = 23, RULE_waitStmt = 24, RULE_assignmentStmt = 25, 
		RULE_argList = 26, RULE_expression = 27, RULE_additiveExpression = 28, 
		RULE_multiplicativeExpression = 29, RULE_primaryExpression = 30;
	private static String[] makeRuleNames() {
		return new String[] {
			"program", "line", "summon", "withAssets", "assetList", "asset", "assetType", 
			"object", "objectBody", "objectSection", "propertiesSection", "propertyBlock", 
			"property", "abilitiesSection", "abilityBlock", "method", "paramList", 
			"parameter", "reactionSection", "type", "statement", "printStmt", "callStmt", 
			"showStmt", "waitStmt", "assignmentStmt", "argList", "expression", "additiveExpression", 
			"multiplicativeExpression", "primaryExpression"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'summon'", "'with'", "'named'", "'as'", "'make'", "'like'", "'has'", 
			"'can'", "'on'", "'is'", "'print'", "'tell'", "'to'", "'show'", "'centered'", 
			"'wait'", "'seconds'", "'and'", "'image'", "'sound'", "'music'", "'video'", 
			"'font'", "'shader'", "'string'", "'number'", "'boolean'", "':'", "'['", 
			"']'", "'('", "')'", "','", "'~'", "'{'", "'}'", "'.'", "'+'", "'-'", 
			"'*'", "'/'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "SUMMON", "WITH", "NAMED", "AS", "MAKE", "LIKE", "HAS", "CAN", 
			"ON", "IS", "PRINT", "TELL", "TO", "SHOW", "CENTERED", "WAIT", "SECONDS", 
			"AND", "IMAGE", "SOUND", "MUSIC", "VIDEO", "FONT", "SHADER", "STRING_TYPE", 
			"NUMBER_TYPE", "BOOLEAN_TYPE", "COLON", "LBRACK", "RBRACK", "LPAREN", 
			"RPAREN", "COMMA", "TILDA", "LBRACE", "RBRACE", "DOT", "PLUS", "MINUS", 
			"MULT", "DIV", "COMMENT", "ID", "NUMBER", "STRING", "NEWLINE", "WS"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "HE_New.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public HE_NewParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ProgramContext extends ParserRuleContext {
		public TerminalNode EOF() { return getToken(HE_NewParser.EOF, 0); }
		public List<LineContext> line() {
			return getRuleContexts(LineContext.class);
		}
		public LineContext line(int i) {
			return getRuleContext(LineContext.class,i);
		}
		public ProgramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_program; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterProgram(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitProgram(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitProgram(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ProgramContext program() throws RecognitionException {
		ProgramContext _localctx = new ProgramContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_program);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(65);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 70368744177702L) != 0)) {
				{
				{
				setState(62);
				line();
				}
				}
				setState(67);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(68);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class LineContext extends ParserRuleContext {
		public SummonContext summon() {
			return getRuleContext(SummonContext.class,0);
		}
		public WithAssetsContext withAssets() {
			return getRuleContext(WithAssetsContext.class,0);
		}
		public ObjectContext object() {
			return getRuleContext(ObjectContext.class,0);
		}
		public TerminalNode NEWLINE() { return getToken(HE_NewParser.NEWLINE, 0); }
		public LineContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_line; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterLine(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitLine(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitLine(this);
			else return visitor.visitChildren(this);
		}
	}

	public final LineContext line() throws RecognitionException {
		LineContext _localctx = new LineContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_line);
		try {
			setState(74);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case SUMMON:
				enterOuterAlt(_localctx, 1);
				{
				setState(70);
				summon();
				}
				break;
			case WITH:
				enterOuterAlt(_localctx, 2);
				{
				setState(71);
				withAssets();
				}
				break;
			case MAKE:
				enterOuterAlt(_localctx, 3);
				{
				setState(72);
				object();
				}
				break;
			case NEWLINE:
				enterOuterAlt(_localctx, 4);
				{
				setState(73);
				match(NEWLINE);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SummonContext extends ParserRuleContext {
		public TerminalNode SUMMON() { return getToken(HE_NewParser.SUMMON, 0); }
		public TerminalNode STRING() { return getToken(HE_NewParser.STRING, 0); }
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode AS() { return getToken(HE_NewParser.AS, 0); }
		public TerminalNode NAMED() { return getToken(HE_NewParser.NAMED, 0); }
		public SummonContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_summon; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterSummon(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitSummon(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitSummon(this);
			else return visitor.visitChildren(this);
		}
	}

	public final SummonContext summon() throws RecognitionException {
		SummonContext _localctx = new SummonContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_summon);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(76);
			match(SUMMON);
			setState(77);
			match(STRING);
			setState(78);
			_la = _input.LA(1);
			if ( !(_la==NAMED || _la==AS) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(79);
			match(ID);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class WithAssetsContext extends ParserRuleContext {
		public TerminalNode WITH() { return getToken(HE_NewParser.WITH, 0); }
		public AssetListContext assetList() {
			return getRuleContext(AssetListContext.class,0);
		}
		public WithAssetsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_withAssets; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterWithAssets(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitWithAssets(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitWithAssets(this);
			else return visitor.visitChildren(this);
		}
	}

	public final WithAssetsContext withAssets() throws RecognitionException {
		WithAssetsContext _localctx = new WithAssetsContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_withAssets);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(81);
			match(WITH);
			setState(82);
			assetList();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AssetListContext extends ParserRuleContext {
		public List<AssetContext> asset() {
			return getRuleContexts(AssetContext.class);
		}
		public AssetContext asset(int i) {
			return getRuleContext(AssetContext.class,i);
		}
		public List<TerminalNode> AND() { return getTokens(HE_NewParser.AND); }
		public TerminalNode AND(int i) {
			return getToken(HE_NewParser.AND, i);
		}
		public AssetListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assetList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAssetList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAssetList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAssetList(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AssetListContext assetList() throws RecognitionException {
		AssetListContext _localctx = new AssetListContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_assetList);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(84);
			asset();
			setState(89);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==AND) {
				{
				{
				setState(85);
				match(AND);
				setState(86);
				asset();
				}
				}
				setState(91);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AssetContext extends ParserRuleContext {
		public AssetTypeContext assetType() {
			return getRuleContext(AssetTypeContext.class,0);
		}
		public TerminalNode STRING() { return getToken(HE_NewParser.STRING, 0); }
		public TerminalNode NAMED() { return getToken(HE_NewParser.NAMED, 0); }
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public AssetContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_asset; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAsset(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAsset(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAsset(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AssetContext asset() throws RecognitionException {
		AssetContext _localctx = new AssetContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_asset);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(92);
			assetType();
			setState(93);
			match(STRING);
			setState(96);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==NAMED) {
				{
				setState(94);
				match(NAMED);
				setState(95);
				match(ID);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AssetTypeContext extends ParserRuleContext {
		public TerminalNode IMAGE() { return getToken(HE_NewParser.IMAGE, 0); }
		public TerminalNode SOUND() { return getToken(HE_NewParser.SOUND, 0); }
		public TerminalNode MUSIC() { return getToken(HE_NewParser.MUSIC, 0); }
		public TerminalNode VIDEO() { return getToken(HE_NewParser.VIDEO, 0); }
		public TerminalNode FONT() { return getToken(HE_NewParser.FONT, 0); }
		public TerminalNode SHADER() { return getToken(HE_NewParser.SHADER, 0); }
		public AssetTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assetType; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAssetType(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAssetType(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAssetType(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AssetTypeContext assetType() throws RecognitionException {
		AssetTypeContext _localctx = new AssetTypeContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_assetType);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(98);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 33030144L) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ObjectContext extends ParserRuleContext {
		public TerminalNode MAKE() { return getToken(HE_NewParser.MAKE, 0); }
		public List<TerminalNode> ID() { return getTokens(HE_NewParser.ID); }
		public TerminalNode ID(int i) {
			return getToken(HE_NewParser.ID, i);
		}
		public TerminalNode LBRACE() { return getToken(HE_NewParser.LBRACE, 0); }
		public ObjectBodyContext objectBody() {
			return getRuleContext(ObjectBodyContext.class,0);
		}
		public TerminalNode RBRACE() { return getToken(HE_NewParser.RBRACE, 0); }
		public TerminalNode LIKE() { return getToken(HE_NewParser.LIKE, 0); }
		public ObjectContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_object; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterObject(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitObject(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitObject(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ObjectContext object() throws RecognitionException {
		ObjectContext _localctx = new ObjectContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_object);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(100);
			match(MAKE);
			setState(101);
			match(ID);
			setState(104);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LIKE) {
				{
				setState(102);
				match(LIKE);
				setState(103);
				match(ID);
				}
			}

			setState(106);
			match(LBRACE);
			setState(107);
			objectBody();
			setState(108);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ObjectBodyContext extends ParserRuleContext {
		public List<ObjectSectionContext> objectSection() {
			return getRuleContexts(ObjectSectionContext.class);
		}
		public ObjectSectionContext objectSection(int i) {
			return getRuleContext(ObjectSectionContext.class,i);
		}
		public ObjectBodyContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_objectBody; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterObjectBody(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitObjectBody(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitObjectBody(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ObjectBodyContext objectBody() throws RecognitionException {
		ObjectBodyContext _localctx = new ObjectBodyContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_objectBody);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(113);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==ON || _la==ID) {
				{
				{
				setState(110);
				objectSection();
				}
				}
				setState(115);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ObjectSectionContext extends ParserRuleContext {
		public PropertiesSectionContext propertiesSection() {
			return getRuleContext(PropertiesSectionContext.class,0);
		}
		public AbilitiesSectionContext abilitiesSection() {
			return getRuleContext(AbilitiesSectionContext.class,0);
		}
		public ReactionSectionContext reactionSection() {
			return getRuleContext(ReactionSectionContext.class,0);
		}
		public ObjectSectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_objectSection; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterObjectSection(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitObjectSection(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitObjectSection(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ObjectSectionContext objectSection() throws RecognitionException {
		ObjectSectionContext _localctx = new ObjectSectionContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_objectSection);
		try {
			setState(119);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,6,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(116);
				propertiesSection();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(117);
				abilitiesSection();
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(118);
				reactionSection();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PropertiesSectionContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode HAS() { return getToken(HE_NewParser.HAS, 0); }
		public TerminalNode COLON() { return getToken(HE_NewParser.COLON, 0); }
		public List<PropertyBlockContext> propertyBlock() {
			return getRuleContexts(PropertyBlockContext.class);
		}
		public PropertyBlockContext propertyBlock(int i) {
			return getRuleContext(PropertyBlockContext.class,i);
		}
		public PropertiesSectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_propertiesSection; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterPropertiesSection(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitPropertiesSection(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitPropertiesSection(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PropertiesSectionContext propertiesSection() throws RecognitionException {
		PropertiesSectionContext _localctx = new PropertiesSectionContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_propertiesSection);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(121);
			match(ID);
			setState(122);
			match(HAS);
			setState(123);
			match(COLON);
			setState(125); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(124);
				propertyBlock();
				}
				}
				setState(127); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( _la==LBRACK );
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PropertyBlockContext extends ParserRuleContext {
		public TerminalNode LBRACK() { return getToken(HE_NewParser.LBRACK, 0); }
		public PropertyContext property() {
			return getRuleContext(PropertyContext.class,0);
		}
		public TerminalNode RBRACK() { return getToken(HE_NewParser.RBRACK, 0); }
		public PropertyBlockContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_propertyBlock; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterPropertyBlock(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitPropertyBlock(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitPropertyBlock(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PropertyBlockContext propertyBlock() throws RecognitionException {
		PropertyBlockContext _localctx = new PropertyBlockContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_propertyBlock);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(129);
			match(LBRACK);
			setState(130);
			property();
			setState(131);
			match(RBRACK);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PropertyContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode IS() { return getToken(HE_NewParser.IS, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public PropertyContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_property; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterProperty(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitProperty(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitProperty(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PropertyContext property() throws RecognitionException {
		PropertyContext _localctx = new PropertyContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_property);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(133);
			match(ID);
			setState(134);
			match(IS);
			setState(135);
			expression();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AbilitiesSectionContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode CAN() { return getToken(HE_NewParser.CAN, 0); }
		public TerminalNode COLON() { return getToken(HE_NewParser.COLON, 0); }
		public List<AbilityBlockContext> abilityBlock() {
			return getRuleContexts(AbilityBlockContext.class);
		}
		public AbilityBlockContext abilityBlock(int i) {
			return getRuleContext(AbilityBlockContext.class,i);
		}
		public AbilitiesSectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_abilitiesSection; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAbilitiesSection(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAbilitiesSection(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAbilitiesSection(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AbilitiesSectionContext abilitiesSection() throws RecognitionException {
		AbilitiesSectionContext _localctx = new AbilitiesSectionContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_abilitiesSection);
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(137);
			match(ID);
			setState(138);
			match(CAN);
			setState(139);
			match(COLON);
			setState(141); 
			_errHandler.sync(this);
			_alt = 1;
			do {
				switch (_alt) {
				case 1:
					{
					{
					setState(140);
					abilityBlock();
					}
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(143); 
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,8,_ctx);
			} while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER );
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AbilityBlockContext extends ParserRuleContext {
		public MethodContext method() {
			return getRuleContext(MethodContext.class,0);
		}
		public TerminalNode LBRACK() { return getToken(HE_NewParser.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(HE_NewParser.RBRACK, 0); }
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public AbilityBlockContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_abilityBlock; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAbilityBlock(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAbilityBlock(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAbilityBlock(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AbilityBlockContext abilityBlock() throws RecognitionException {
		AbilityBlockContext _localctx = new AbilityBlockContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_abilityBlock);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(145);
			method();
			setState(146);
			match(LBRACK);
			setState(150);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 8796093110272L) != 0)) {
				{
				{
				setState(147);
				statement();
				}
				}
				setState(152);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(153);
			match(RBRACK);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class MethodContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode LPAREN() { return getToken(HE_NewParser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(HE_NewParser.RPAREN, 0); }
		public ParamListContext paramList() {
			return getRuleContext(ParamListContext.class,0);
		}
		public MethodContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_method; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterMethod(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitMethod(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitMethod(this);
			else return visitor.visitChildren(this);
		}
	}

	public final MethodContext method() throws RecognitionException {
		MethodContext _localctx = new MethodContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_method);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(155);
			match(ID);
			setState(156);
			match(LPAREN);
			setState(158);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ID) {
				{
				setState(157);
				paramList();
				}
			}

			setState(160);
			match(RPAREN);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ParamListContext extends ParserRuleContext {
		public List<ParameterContext> parameter() {
			return getRuleContexts(ParameterContext.class);
		}
		public ParameterContext parameter(int i) {
			return getRuleContext(ParameterContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(HE_NewParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(HE_NewParser.COMMA, i);
		}
		public ParamListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_paramList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterParamList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitParamList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitParamList(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ParamListContext paramList() throws RecognitionException {
		ParamListContext _localctx = new ParamListContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_paramList);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(162);
			parameter();
			setState(167);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(163);
				match(COMMA);
				setState(164);
				parameter();
				}
				}
				setState(169);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ParameterContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode COLON() { return getToken(HE_NewParser.COLON, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public ParameterContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_parameter; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterParameter(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitParameter(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitParameter(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ParameterContext parameter() throws RecognitionException {
		ParameterContext _localctx = new ParameterContext(_ctx, getState());
		enterRule(_localctx, 34, RULE_parameter);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(170);
			match(ID);
			setState(171);
			match(COLON);
			setState(172);
			type();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ReactionSectionContext extends ParserRuleContext {
		public TerminalNode ON() { return getToken(HE_NewParser.ON, 0); }
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode LBRACK() { return getToken(HE_NewParser.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(HE_NewParser.RBRACK, 0); }
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public ReactionSectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_reactionSection; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterReactionSection(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitReactionSection(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitReactionSection(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ReactionSectionContext reactionSection() throws RecognitionException {
		ReactionSectionContext _localctx = new ReactionSectionContext(_ctx, getState());
		enterRule(_localctx, 36, RULE_reactionSection);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(174);
			match(ON);
			setState(175);
			match(ID);
			setState(176);
			match(LBRACK);
			setState(180);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 8796093110272L) != 0)) {
				{
				{
				setState(177);
				statement();
				}
				}
				setState(182);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(183);
			match(RBRACK);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class TypeContext extends ParserRuleContext {
		public TerminalNode STRING_TYPE() { return getToken(HE_NewParser.STRING_TYPE, 0); }
		public TerminalNode NUMBER_TYPE() { return getToken(HE_NewParser.NUMBER_TYPE, 0); }
		public TerminalNode BOOLEAN_TYPE() { return getToken(HE_NewParser.BOOLEAN_TYPE, 0); }
		public TypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_type; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterType(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitType(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitType(this);
			else return visitor.visitChildren(this);
		}
	}

	public final TypeContext type() throws RecognitionException {
		TypeContext _localctx = new TypeContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_type);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(185);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 234881024L) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class StatementContext extends ParserRuleContext {
		public PrintStmtContext printStmt() {
			return getRuleContext(PrintStmtContext.class,0);
		}
		public CallStmtContext callStmt() {
			return getRuleContext(CallStmtContext.class,0);
		}
		public ShowStmtContext showStmt() {
			return getRuleContext(ShowStmtContext.class,0);
		}
		public WaitStmtContext waitStmt() {
			return getRuleContext(WaitStmtContext.class,0);
		}
		public AssignmentStmtContext assignmentStmt() {
			return getRuleContext(AssignmentStmtContext.class,0);
		}
		public StatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_statement; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterStatement(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitStatement(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitStatement(this);
			else return visitor.visitChildren(this);
		}
	}

	public final StatementContext statement() throws RecognitionException {
		StatementContext _localctx = new StatementContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_statement);
		try {
			setState(192);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case PRINT:
				enterOuterAlt(_localctx, 1);
				{
				setState(187);
				printStmt();
				}
				break;
			case TELL:
				enterOuterAlt(_localctx, 2);
				{
				setState(188);
				callStmt();
				}
				break;
			case SHOW:
				enterOuterAlt(_localctx, 3);
				{
				setState(189);
				showStmt();
				}
				break;
			case WAIT:
				enterOuterAlt(_localctx, 4);
				{
				setState(190);
				waitStmt();
				}
				break;
			case ID:
				enterOuterAlt(_localctx, 5);
				{
				setState(191);
				assignmentStmt();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PrintStmtContext extends ParserRuleContext {
		public TerminalNode PRINT() { return getToken(HE_NewParser.PRINT, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public PrintStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_printStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterPrintStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitPrintStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitPrintStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PrintStmtContext printStmt() throws RecognitionException {
		PrintStmtContext _localctx = new PrintStmtContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_printStmt);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(194);
			match(PRINT);
			setState(195);
			expression();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class CallStmtContext extends ParserRuleContext {
		public TerminalNode TELL() { return getToken(HE_NewParser.TELL, 0); }
		public List<TerminalNode> ID() { return getTokens(HE_NewParser.ID); }
		public TerminalNode ID(int i) {
			return getToken(HE_NewParser.ID, i);
		}
		public TerminalNode TO() { return getToken(HE_NewParser.TO, 0); }
		public TerminalNode WITH() { return getToken(HE_NewParser.WITH, 0); }
		public ArgListContext argList() {
			return getRuleContext(ArgListContext.class,0);
		}
		public CallStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_callStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterCallStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitCallStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitCallStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final CallStmtContext callStmt() throws RecognitionException {
		CallStmtContext _localctx = new CallStmtContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_callStmt);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(197);
			match(TELL);
			setState(198);
			match(ID);
			setState(199);
			match(TO);
			setState(200);
			match(ID);
			setState(203);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==WITH) {
				{
				setState(201);
				match(WITH);
				setState(202);
				argList();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ShowStmtContext extends ParserRuleContext {
		public TerminalNode SHOW() { return getToken(HE_NewParser.SHOW, 0); }
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode CENTERED() { return getToken(HE_NewParser.CENTERED, 0); }
		public ShowStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_showStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterShowStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitShowStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitShowStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ShowStmtContext showStmt() throws RecognitionException {
		ShowStmtContext _localctx = new ShowStmtContext(_ctx, getState());
		enterRule(_localctx, 46, RULE_showStmt);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(205);
			match(SHOW);
			setState(206);
			match(ID);
			setState(208);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==CENTERED) {
				{
				setState(207);
				match(CENTERED);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class WaitStmtContext extends ParserRuleContext {
		public TerminalNode WAIT() { return getToken(HE_NewParser.WAIT, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SECONDS() { return getToken(HE_NewParser.SECONDS, 0); }
		public WaitStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_waitStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterWaitStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitWaitStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitWaitStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final WaitStmtContext waitStmt() throws RecognitionException {
		WaitStmtContext _localctx = new WaitStmtContext(_ctx, getState());
		enterRule(_localctx, 48, RULE_waitStmt);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(210);
			match(WAIT);
			setState(211);
			expression();
			setState(212);
			match(SECONDS);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AssignmentStmtContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode IS() { return getToken(HE_NewParser.IS, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public AssignmentStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assignmentStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAssignmentStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAssignmentStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAssignmentStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AssignmentStmtContext assignmentStmt() throws RecognitionException {
		AssignmentStmtContext _localctx = new AssignmentStmtContext(_ctx, getState());
		enterRule(_localctx, 50, RULE_assignmentStmt);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(214);
			match(ID);
			setState(215);
			match(IS);
			setState(216);
			expression();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArgListContext extends ParserRuleContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(HE_NewParser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(HE_NewParser.COMMA, i);
		}
		public ArgListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_argList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterArgList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitArgList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitArgList(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ArgListContext argList() throws RecognitionException {
		ArgListContext _localctx = new ArgListContext(_ctx, getState());
		enterRule(_localctx, 52, RULE_argList);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(218);
			expression();
			setState(223);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(219);
				match(COMMA);
				setState(220);
				expression();
				}
				}
				setState(225);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExpressionContext extends ParserRuleContext {
		public AdditiveExpressionContext additiveExpression() {
			return getRuleContext(AdditiveExpressionContext.class,0);
		}
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ExpressionContext expression() throws RecognitionException {
		ExpressionContext _localctx = new ExpressionContext(_ctx, getState());
		enterRule(_localctx, 54, RULE_expression);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(226);
			additiveExpression();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AdditiveExpressionContext extends ParserRuleContext {
		public List<MultiplicativeExpressionContext> multiplicativeExpression() {
			return getRuleContexts(MultiplicativeExpressionContext.class);
		}
		public MultiplicativeExpressionContext multiplicativeExpression(int i) {
			return getRuleContext(MultiplicativeExpressionContext.class,i);
		}
		public List<TerminalNode> PLUS() { return getTokens(HE_NewParser.PLUS); }
		public TerminalNode PLUS(int i) {
			return getToken(HE_NewParser.PLUS, i);
		}
		public List<TerminalNode> MINUS() { return getTokens(HE_NewParser.MINUS); }
		public TerminalNode MINUS(int i) {
			return getToken(HE_NewParser.MINUS, i);
		}
		public AdditiveExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_additiveExpression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterAdditiveExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitAdditiveExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitAdditiveExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AdditiveExpressionContext additiveExpression() throws RecognitionException {
		AdditiveExpressionContext _localctx = new AdditiveExpressionContext(_ctx, getState());
		enterRule(_localctx, 56, RULE_additiveExpression);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(228);
			multiplicativeExpression();
			setState(233);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==PLUS || _la==MINUS) {
				{
				{
				setState(229);
				_la = _input.LA(1);
				if ( !(_la==PLUS || _la==MINUS) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(230);
				multiplicativeExpression();
				}
				}
				setState(235);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class MultiplicativeExpressionContext extends ParserRuleContext {
		public List<PrimaryExpressionContext> primaryExpression() {
			return getRuleContexts(PrimaryExpressionContext.class);
		}
		public PrimaryExpressionContext primaryExpression(int i) {
			return getRuleContext(PrimaryExpressionContext.class,i);
		}
		public List<TerminalNode> MULT() { return getTokens(HE_NewParser.MULT); }
		public TerminalNode MULT(int i) {
			return getToken(HE_NewParser.MULT, i);
		}
		public List<TerminalNode> DIV() { return getTokens(HE_NewParser.DIV); }
		public TerminalNode DIV(int i) {
			return getToken(HE_NewParser.DIV, i);
		}
		public MultiplicativeExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_multiplicativeExpression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterMultiplicativeExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitMultiplicativeExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitMultiplicativeExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final MultiplicativeExpressionContext multiplicativeExpression() throws RecognitionException {
		MultiplicativeExpressionContext _localctx = new MultiplicativeExpressionContext(_ctx, getState());
		enterRule(_localctx, 58, RULE_multiplicativeExpression);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(236);
			primaryExpression();
			setState(241);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==MULT || _la==DIV) {
				{
				{
				setState(237);
				_la = _input.LA(1);
				if ( !(_la==MULT || _la==DIV) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(238);
				primaryExpression();
				}
				}
				setState(243);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PrimaryExpressionContext extends ParserRuleContext {
		public TerminalNode STRING() { return getToken(HE_NewParser.STRING, 0); }
		public TerminalNode ID() { return getToken(HE_NewParser.ID, 0); }
		public TerminalNode NUMBER() { return getToken(HE_NewParser.NUMBER, 0); }
		public TerminalNode LPAREN() { return getToken(HE_NewParser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(HE_NewParser.RPAREN, 0); }
		public PrimaryExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_primaryExpression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).enterPrimaryExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HE_NewListener ) ((HE_NewListener)listener).exitPrimaryExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HE_NewVisitor ) return ((HE_NewVisitor<? extends T>)visitor).visitPrimaryExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PrimaryExpressionContext primaryExpression() throws RecognitionException {
		PrimaryExpressionContext _localctx = new PrimaryExpressionContext(_ctx, getState());
		enterRule(_localctx, 60, RULE_primaryExpression);
		try {
			setState(251);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case STRING:
				enterOuterAlt(_localctx, 1);
				{
				setState(244);
				match(STRING);
				}
				break;
			case ID:
				enterOuterAlt(_localctx, 2);
				{
				setState(245);
				match(ID);
				}
				break;
			case NUMBER:
				enterOuterAlt(_localctx, 3);
				{
				setState(246);
				match(NUMBER);
				}
				break;
			case LPAREN:
				enterOuterAlt(_localctx, 4);
				{
				setState(247);
				match(LPAREN);
				setState(248);
				expression();
				setState(249);
				match(RPAREN);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static final String _serializedATN =
		"\u0004\u0001/\u00fe\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001\u0002"+
		"\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004\u0002"+
		"\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007\u0002"+
		"\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b\u0002"+
		"\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007\u000f"+
		"\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0002\u0012\u0007\u0012"+
		"\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014\u0002\u0015\u0007\u0015"+
		"\u0002\u0016\u0007\u0016\u0002\u0017\u0007\u0017\u0002\u0018\u0007\u0018"+
		"\u0002\u0019\u0007\u0019\u0002\u001a\u0007\u001a\u0002\u001b\u0007\u001b"+
		"\u0002\u001c\u0007\u001c\u0002\u001d\u0007\u001d\u0002\u001e\u0007\u001e"+
		"\u0001\u0000\u0005\u0000@\b\u0000\n\u0000\f\u0000C\t\u0000\u0001\u0000"+
		"\u0001\u0000\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0003\u0001"+
		"K\b\u0001\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0003\u0001\u0003\u0001\u0003\u0001\u0004\u0001\u0004\u0001\u0004"+
		"\u0005\u0004X\b\u0004\n\u0004\f\u0004[\t\u0004\u0001\u0005\u0001\u0005"+
		"\u0001\u0005\u0001\u0005\u0003\u0005a\b\u0005\u0001\u0006\u0001\u0006"+
		"\u0001\u0007\u0001\u0007\u0001\u0007\u0001\u0007\u0003\u0007i\b\u0007"+
		"\u0001\u0007\u0001\u0007\u0001\u0007\u0001\u0007\u0001\b\u0005\bp\b\b"+
		"\n\b\f\bs\t\b\u0001\t\u0001\t\u0001\t\u0003\tx\b\t\u0001\n\u0001\n\u0001"+
		"\n\u0001\n\u0004\n~\b\n\u000b\n\f\n\u007f\u0001\u000b\u0001\u000b\u0001"+
		"\u000b\u0001\u000b\u0001\f\u0001\f\u0001\f\u0001\f\u0001\r\u0001\r\u0001"+
		"\r\u0001\r\u0004\r\u008e\b\r\u000b\r\f\r\u008f\u0001\u000e\u0001\u000e"+
		"\u0001\u000e\u0005\u000e\u0095\b\u000e\n\u000e\f\u000e\u0098\t\u000e\u0001"+
		"\u000e\u0001\u000e\u0001\u000f\u0001\u000f\u0001\u000f\u0003\u000f\u009f"+
		"\b\u000f\u0001\u000f\u0001\u000f\u0001\u0010\u0001\u0010\u0001\u0010\u0005"+
		"\u0010\u00a6\b\u0010\n\u0010\f\u0010\u00a9\t\u0010\u0001\u0011\u0001\u0011"+
		"\u0001\u0011\u0001\u0011\u0001\u0012\u0001\u0012\u0001\u0012\u0001\u0012"+
		"\u0005\u0012\u00b3\b\u0012\n\u0012\f\u0012\u00b6\t\u0012\u0001\u0012\u0001"+
		"\u0012\u0001\u0013\u0001\u0013\u0001\u0014\u0001\u0014\u0001\u0014\u0001"+
		"\u0014\u0001\u0014\u0003\u0014\u00c1\b\u0014\u0001\u0015\u0001\u0015\u0001"+
		"\u0015\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001"+
		"\u0016\u0003\u0016\u00cc\b\u0016\u0001\u0017\u0001\u0017\u0001\u0017\u0003"+
		"\u0017\u00d1\b\u0017\u0001\u0018\u0001\u0018\u0001\u0018\u0001\u0018\u0001"+
		"\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u001a\u0001\u001a\u0001"+
		"\u001a\u0005\u001a\u00de\b\u001a\n\u001a\f\u001a\u00e1\t\u001a\u0001\u001b"+
		"\u0001\u001b\u0001\u001c\u0001\u001c\u0001\u001c\u0005\u001c\u00e8\b\u001c"+
		"\n\u001c\f\u001c\u00eb\t\u001c\u0001\u001d\u0001\u001d\u0001\u001d\u0005"+
		"\u001d\u00f0\b\u001d\n\u001d\f\u001d\u00f3\t\u001d\u0001\u001e\u0001\u001e"+
		"\u0001\u001e\u0001\u001e\u0001\u001e\u0001\u001e\u0001\u001e\u0003\u001e"+
		"\u00fc\b\u001e\u0001\u001e\u0000\u0000\u001f\u0000\u0002\u0004\u0006\b"+
		"\n\f\u000e\u0010\u0012\u0014\u0016\u0018\u001a\u001c\u001e \"$&(*,.02"+
		"468:<\u0000\u0005\u0001\u0000\u0003\u0004\u0001\u0000\u0013\u0018\u0001"+
		"\u0000\u0019\u001b\u0001\u0000&\'\u0001\u0000()\u00fa\u0000A\u0001\u0000"+
		"\u0000\u0000\u0002J\u0001\u0000\u0000\u0000\u0004L\u0001\u0000\u0000\u0000"+
		"\u0006Q\u0001\u0000\u0000\u0000\bT\u0001\u0000\u0000\u0000\n\\\u0001\u0000"+
		"\u0000\u0000\fb\u0001\u0000\u0000\u0000\u000ed\u0001\u0000\u0000\u0000"+
		"\u0010q\u0001\u0000\u0000\u0000\u0012w\u0001\u0000\u0000\u0000\u0014y"+
		"\u0001\u0000\u0000\u0000\u0016\u0081\u0001\u0000\u0000\u0000\u0018\u0085"+
		"\u0001\u0000\u0000\u0000\u001a\u0089\u0001\u0000\u0000\u0000\u001c\u0091"+
		"\u0001\u0000\u0000\u0000\u001e\u009b\u0001\u0000\u0000\u0000 \u00a2\u0001"+
		"\u0000\u0000\u0000\"\u00aa\u0001\u0000\u0000\u0000$\u00ae\u0001\u0000"+
		"\u0000\u0000&\u00b9\u0001\u0000\u0000\u0000(\u00c0\u0001\u0000\u0000\u0000"+
		"*\u00c2\u0001\u0000\u0000\u0000,\u00c5\u0001\u0000\u0000\u0000.\u00cd"+
		"\u0001\u0000\u0000\u00000\u00d2\u0001\u0000\u0000\u00002\u00d6\u0001\u0000"+
		"\u0000\u00004\u00da\u0001\u0000\u0000\u00006\u00e2\u0001\u0000\u0000\u0000"+
		"8\u00e4\u0001\u0000\u0000\u0000:\u00ec\u0001\u0000\u0000\u0000<\u00fb"+
		"\u0001\u0000\u0000\u0000>@\u0003\u0002\u0001\u0000?>\u0001\u0000\u0000"+
		"\u0000@C\u0001\u0000\u0000\u0000A?\u0001\u0000\u0000\u0000AB\u0001\u0000"+
		"\u0000\u0000BD\u0001\u0000\u0000\u0000CA\u0001\u0000\u0000\u0000DE\u0005"+
		"\u0000\u0000\u0001E\u0001\u0001\u0000\u0000\u0000FK\u0003\u0004\u0002"+
		"\u0000GK\u0003\u0006\u0003\u0000HK\u0003\u000e\u0007\u0000IK\u0005.\u0000"+
		"\u0000JF\u0001\u0000\u0000\u0000JG\u0001\u0000\u0000\u0000JH\u0001\u0000"+
		"\u0000\u0000JI\u0001\u0000\u0000\u0000K\u0003\u0001\u0000\u0000\u0000"+
		"LM\u0005\u0001\u0000\u0000MN\u0005-\u0000\u0000NO\u0007\u0000\u0000\u0000"+
		"OP\u0005+\u0000\u0000P\u0005\u0001\u0000\u0000\u0000QR\u0005\u0002\u0000"+
		"\u0000RS\u0003\b\u0004\u0000S\u0007\u0001\u0000\u0000\u0000TY\u0003\n"+
		"\u0005\u0000UV\u0005\u0012\u0000\u0000VX\u0003\n\u0005\u0000WU\u0001\u0000"+
		"\u0000\u0000X[\u0001\u0000\u0000\u0000YW\u0001\u0000\u0000\u0000YZ\u0001"+
		"\u0000\u0000\u0000Z\t\u0001\u0000\u0000\u0000[Y\u0001\u0000\u0000\u0000"+
		"\\]\u0003\f\u0006\u0000]`\u0005-\u0000\u0000^_\u0005\u0003\u0000\u0000"+
		"_a\u0005+\u0000\u0000`^\u0001\u0000\u0000\u0000`a\u0001\u0000\u0000\u0000"+
		"a\u000b\u0001\u0000\u0000\u0000bc\u0007\u0001\u0000\u0000c\r\u0001\u0000"+
		"\u0000\u0000de\u0005\u0005\u0000\u0000eh\u0005+\u0000\u0000fg\u0005\u0006"+
		"\u0000\u0000gi\u0005+\u0000\u0000hf\u0001\u0000\u0000\u0000hi\u0001\u0000"+
		"\u0000\u0000ij\u0001\u0000\u0000\u0000jk\u0005#\u0000\u0000kl\u0003\u0010"+
		"\b\u0000lm\u0005$\u0000\u0000m\u000f\u0001\u0000\u0000\u0000np\u0003\u0012"+
		"\t\u0000on\u0001\u0000\u0000\u0000ps\u0001\u0000\u0000\u0000qo\u0001\u0000"+
		"\u0000\u0000qr\u0001\u0000\u0000\u0000r\u0011\u0001\u0000\u0000\u0000"+
		"sq\u0001\u0000\u0000\u0000tx\u0003\u0014\n\u0000ux\u0003\u001a\r\u0000"+
		"vx\u0003$\u0012\u0000wt\u0001\u0000\u0000\u0000wu\u0001\u0000\u0000\u0000"+
		"wv\u0001\u0000\u0000\u0000x\u0013\u0001\u0000\u0000\u0000yz\u0005+\u0000"+
		"\u0000z{\u0005\u0007\u0000\u0000{}\u0005\u001c\u0000\u0000|~\u0003\u0016"+
		"\u000b\u0000}|\u0001\u0000\u0000\u0000~\u007f\u0001\u0000\u0000\u0000"+
		"\u007f}\u0001\u0000\u0000\u0000\u007f\u0080\u0001\u0000\u0000\u0000\u0080"+
		"\u0015\u0001\u0000\u0000\u0000\u0081\u0082\u0005\u001d\u0000\u0000\u0082"+
		"\u0083\u0003\u0018\f\u0000\u0083\u0084\u0005\u001e\u0000\u0000\u0084\u0017"+
		"\u0001\u0000\u0000\u0000\u0085\u0086\u0005+\u0000\u0000\u0086\u0087\u0005"+
		"\n\u0000\u0000\u0087\u0088\u00036\u001b\u0000\u0088\u0019\u0001\u0000"+
		"\u0000\u0000\u0089\u008a\u0005+\u0000\u0000\u008a\u008b\u0005\b\u0000"+
		"\u0000\u008b\u008d\u0005\u001c\u0000\u0000\u008c\u008e\u0003\u001c\u000e"+
		"\u0000\u008d\u008c\u0001\u0000\u0000\u0000\u008e\u008f\u0001\u0000\u0000"+
		"\u0000\u008f\u008d\u0001\u0000\u0000\u0000\u008f\u0090\u0001\u0000\u0000"+
		"\u0000\u0090\u001b\u0001\u0000\u0000\u0000\u0091\u0092\u0003\u001e\u000f"+
		"\u0000\u0092\u0096\u0005\u001d\u0000\u0000\u0093\u0095\u0003(\u0014\u0000"+
		"\u0094\u0093\u0001\u0000\u0000\u0000\u0095\u0098\u0001\u0000\u0000\u0000"+
		"\u0096\u0094\u0001\u0000\u0000\u0000\u0096\u0097\u0001\u0000\u0000\u0000"+
		"\u0097\u0099\u0001\u0000\u0000\u0000\u0098\u0096\u0001\u0000\u0000\u0000"+
		"\u0099\u009a\u0005\u001e\u0000\u0000\u009a\u001d\u0001\u0000\u0000\u0000"+
		"\u009b\u009c\u0005+\u0000\u0000\u009c\u009e\u0005\u001f\u0000\u0000\u009d"+
		"\u009f\u0003 \u0010\u0000\u009e\u009d\u0001\u0000\u0000\u0000\u009e\u009f"+
		"\u0001\u0000\u0000\u0000\u009f\u00a0\u0001\u0000\u0000\u0000\u00a0\u00a1"+
		"\u0005 \u0000\u0000\u00a1\u001f\u0001\u0000\u0000\u0000\u00a2\u00a7\u0003"+
		"\"\u0011\u0000\u00a3\u00a4\u0005!\u0000\u0000\u00a4\u00a6\u0003\"\u0011"+
		"\u0000\u00a5\u00a3\u0001\u0000\u0000\u0000\u00a6\u00a9\u0001\u0000\u0000"+
		"\u0000\u00a7\u00a5\u0001\u0000\u0000\u0000\u00a7\u00a8\u0001\u0000\u0000"+
		"\u0000\u00a8!\u0001\u0000\u0000\u0000\u00a9\u00a7\u0001\u0000\u0000\u0000"+
		"\u00aa\u00ab\u0005+\u0000\u0000\u00ab\u00ac\u0005\u001c\u0000\u0000\u00ac"+
		"\u00ad\u0003&\u0013\u0000\u00ad#\u0001\u0000\u0000\u0000\u00ae\u00af\u0005"+
		"\t\u0000\u0000\u00af\u00b0\u0005+\u0000\u0000\u00b0\u00b4\u0005\u001d"+
		"\u0000\u0000\u00b1\u00b3\u0003(\u0014\u0000\u00b2\u00b1\u0001\u0000\u0000"+
		"\u0000\u00b3\u00b6\u0001\u0000\u0000\u0000\u00b4\u00b2\u0001\u0000\u0000"+
		"\u0000\u00b4\u00b5\u0001\u0000\u0000\u0000\u00b5\u00b7\u0001\u0000\u0000"+
		"\u0000\u00b6\u00b4\u0001\u0000\u0000\u0000\u00b7\u00b8\u0005\u001e\u0000"+
		"\u0000\u00b8%\u0001\u0000\u0000\u0000\u00b9\u00ba\u0007\u0002\u0000\u0000"+
		"\u00ba\'\u0001\u0000\u0000\u0000\u00bb\u00c1\u0003*\u0015\u0000\u00bc"+
		"\u00c1\u0003,\u0016\u0000\u00bd\u00c1\u0003.\u0017\u0000\u00be\u00c1\u0003"+
		"0\u0018\u0000\u00bf\u00c1\u00032\u0019\u0000\u00c0\u00bb\u0001\u0000\u0000"+
		"\u0000\u00c0\u00bc\u0001\u0000\u0000\u0000\u00c0\u00bd\u0001\u0000\u0000"+
		"\u0000\u00c0\u00be\u0001\u0000\u0000\u0000\u00c0\u00bf\u0001\u0000\u0000"+
		"\u0000\u00c1)\u0001\u0000\u0000\u0000\u00c2\u00c3\u0005\u000b\u0000\u0000"+
		"\u00c3\u00c4\u00036\u001b\u0000\u00c4+\u0001\u0000\u0000\u0000\u00c5\u00c6"+
		"\u0005\f\u0000\u0000\u00c6\u00c7\u0005+\u0000\u0000\u00c7\u00c8\u0005"+
		"\r\u0000\u0000\u00c8\u00cb\u0005+\u0000\u0000\u00c9\u00ca\u0005\u0002"+
		"\u0000\u0000\u00ca\u00cc\u00034\u001a\u0000\u00cb\u00c9\u0001\u0000\u0000"+
		"\u0000\u00cb\u00cc\u0001\u0000\u0000\u0000\u00cc-\u0001\u0000\u0000\u0000"+
		"\u00cd\u00ce\u0005\u000e\u0000\u0000\u00ce\u00d0\u0005+\u0000\u0000\u00cf"+
		"\u00d1\u0005\u000f\u0000\u0000\u00d0\u00cf\u0001\u0000\u0000\u0000\u00d0"+
		"\u00d1\u0001\u0000\u0000\u0000\u00d1/\u0001\u0000\u0000\u0000\u00d2\u00d3"+
		"\u0005\u0010\u0000\u0000\u00d3\u00d4\u00036\u001b\u0000\u00d4\u00d5\u0005"+
		"\u0011\u0000\u0000\u00d51\u0001\u0000\u0000\u0000\u00d6\u00d7\u0005+\u0000"+
		"\u0000\u00d7\u00d8\u0005\n\u0000\u0000\u00d8\u00d9\u00036\u001b\u0000"+
		"\u00d93\u0001\u0000\u0000\u0000\u00da\u00df\u00036\u001b\u0000\u00db\u00dc"+
		"\u0005!\u0000\u0000\u00dc\u00de\u00036\u001b\u0000\u00dd\u00db\u0001\u0000"+
		"\u0000\u0000\u00de\u00e1\u0001\u0000\u0000\u0000\u00df\u00dd\u0001\u0000"+
		"\u0000\u0000\u00df\u00e0\u0001\u0000\u0000\u0000\u00e05\u0001\u0000\u0000"+
		"\u0000\u00e1\u00df\u0001\u0000\u0000\u0000\u00e2\u00e3\u00038\u001c\u0000"+
		"\u00e37\u0001\u0000\u0000\u0000\u00e4\u00e9\u0003:\u001d\u0000\u00e5\u00e6"+
		"\u0007\u0003\u0000\u0000\u00e6\u00e8\u0003:\u001d\u0000\u00e7\u00e5\u0001"+
		"\u0000\u0000\u0000\u00e8\u00eb\u0001\u0000\u0000\u0000\u00e9\u00e7\u0001"+
		"\u0000\u0000\u0000\u00e9\u00ea\u0001\u0000\u0000\u0000\u00ea9\u0001\u0000"+
		"\u0000\u0000\u00eb\u00e9\u0001\u0000\u0000\u0000\u00ec\u00f1\u0003<\u001e"+
		"\u0000\u00ed\u00ee\u0007\u0004\u0000\u0000\u00ee\u00f0\u0003<\u001e\u0000"+
		"\u00ef\u00ed\u0001\u0000\u0000\u0000\u00f0\u00f3\u0001\u0000\u0000\u0000"+
		"\u00f1\u00ef\u0001\u0000\u0000\u0000\u00f1\u00f2\u0001\u0000\u0000\u0000"+
		"\u00f2;\u0001\u0000\u0000\u0000\u00f3\u00f1\u0001\u0000\u0000\u0000\u00f4"+
		"\u00fc\u0005-\u0000\u0000\u00f5\u00fc\u0005+\u0000\u0000\u00f6\u00fc\u0005"+
		",\u0000\u0000\u00f7\u00f8\u0005\u001f\u0000\u0000\u00f8\u00f9\u00036\u001b"+
		"\u0000\u00f9\u00fa\u0005 \u0000\u0000\u00fa\u00fc\u0001\u0000\u0000\u0000"+
		"\u00fb\u00f4\u0001\u0000\u0000\u0000\u00fb\u00f5\u0001\u0000\u0000\u0000"+
		"\u00fb\u00f6\u0001\u0000\u0000\u0000\u00fb\u00f7\u0001\u0000\u0000\u0000"+
		"\u00fc=\u0001\u0000\u0000\u0000\u0014AJY`hqw\u007f\u008f\u0096\u009e\u00a7"+
		"\u00b4\u00c0\u00cb\u00d0\u00df\u00e9\u00f1\u00fb";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}