// Generated from HE.g4 by ANTLR 4.13.0
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class HEParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.0", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, T__4=5, T__5=6, T__6=7, T__7=8, T__8=9, 
		T__9=10, T__10=11, T__11=12, T__12=13, T__13=14, T__14=15, T__15=16, T__16=17, 
		ID=18, NUMBER=19, STRING=20, BOOLEAN=21, SUMMON=22, WITH=23, NAMED=24, 
		AS=25, CREATE=26, MAKE=27, LIKE=28, HAS=29, OWNS=30, CARRIES=31, CAN=32, 
		KNOWS_HOW_TO=33, WHEN=34, WHENEVER=35, ON=36, REMEMBERS=37, IS=38, STARTS_AS=39, 
		RETURNS=40, SET=41, TO=42, IF=43, THEN=44, ELSE=45, REPEAT=46, WHILE=47, 
		TELL=48, WAIT=49, RETURN=50, OR=51, AND=52, NOT=53, IMAGE=54, SOUND=55, 
		MUSIC=56, VIDEO=57, FONT=58, SHADER=59, SAY=60, PRINT=61, SECONDS=62, 
		FRAMES=63, NUMBER_TYPE=64, STRING_TYPE=65, BOOLEAN_TYPE=66, PLUS=67, MINUS=68, 
		STAR=69, SLASH=70, BANG=71, NEWLINE=72, WS=73;
	public static final int
		RULE_program = 0, RULE_line = 1, RULE_summon = 2, RULE_withAssets = 3, 
		RULE_assetList = 4, RULE_asset = 5, RULE_assetType = 6, RULE_object = 7, 
		RULE_objectBody = 8, RULE_section = 9, RULE_properties = 10, RULE_abilities = 11, 
		RULE_reactions = 12, RULE_memories = 13, RULE_trigger = 14, RULE_block = 15, 
		RULE_property = 16, RULE_action = 17, RULE_paramList = 18, RULE_parameter = 19, 
		RULE_statement = 20, RULE_say = 21, RULE_change = 22, RULE_decide = 23, 
		RULE_repeat = 24, RULE_call = 25, RULE_waitStmt = 26, RULE_returnStmt = 27, 
		RULE_argList = 28, RULE_expression = 29, RULE_logicOr = 30, RULE_logicAnd = 31, 
		RULE_comparison = 32, RULE_arithmetic = 33, RULE_term = 34, RULE_factor = 35, 
		RULE_unary = 36, RULE_primary = 37, RULE_callExpression = 38, RULE_arrayLiteral = 39, 
		RULE_type = 40, RULE_baseType = 41, RULE_globalStatement = 42, RULE_comment = 43, 
		RULE_compOp = 44, RULE_addOp = 45, RULE_multOp = 46;
	private static String[] makeRuleNames() {
		return new String[] {
			"program", "line", "summon", "withAssets", "assetList", "asset", "assetType", 
			"object", "objectBody", "section", "properties", "abilities", "reactions", 
			"memories", "trigger", "block", "property", "action", "paramList", "parameter", 
			"statement", "say", "change", "decide", "repeat", "call", "waitStmt", 
			"returnStmt", "argList", "expression", "logicOr", "logicAnd", "comparison", 
			"arithmetic", "term", "factor", "unary", "primary", "callExpression", 
			"arrayLiteral", "type", "baseType", "globalStatement", "comment", "compOp", 
			"addOp", "multOp"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'{'", "'}'", "'['", "']'", "':'", "'('", "')'", "','", "'^'", 
			"'[]'", "'~'", "'=='", "'!='", "'>'", "'<'", "'>='", "'<='", null, null, 
			null, null, "'summon'", "'with'", "'named'", "'as'", "'create'", "'make'", 
			"'like'", "'has'", "'owns'", "'carries'", "'can'", null, "'when'", "'whenever'", 
			"'on'", "'remembers'", "'is'", "'starts as'", "'returns'", "'set'", "'to'", 
			"'if'", "'then'", "'else'", "'repeat'", "'while'", "'tell'", "'wait'", 
			"'return'", "'or'", "'and'", "'not'", "'image'", "'sound'", "'music'", 
			"'video'", "'font'", "'shader'", "'say'", "'print'", "'seconds'", "'frames'", 
			"'number'", "'string'", "'boolean'", "'+'", "'-'", "'*'", "'/'", "'!'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, null, null, null, null, null, null, null, null, null, null, null, 
			null, null, null, null, null, null, "ID", "NUMBER", "STRING", "BOOLEAN", 
			"SUMMON", "WITH", "NAMED", "AS", "CREATE", "MAKE", "LIKE", "HAS", "OWNS", 
			"CARRIES", "CAN", "KNOWS_HOW_TO", "WHEN", "WHENEVER", "ON", "REMEMBERS", 
			"IS", "STARTS_AS", "RETURNS", "SET", "TO", "IF", "THEN", "ELSE", "REPEAT", 
			"WHILE", "TELL", "WAIT", "RETURN", "OR", "AND", "NOT", "IMAGE", "SOUND", 
			"MUSIC", "VIDEO", "FONT", "SHADER", "SAY", "PRINT", "SECONDS", "FRAMES", 
			"NUMBER_TYPE", "STRING_TYPE", "BOOLEAN_TYPE", "PLUS", "MINUS", "STAR", 
			"SLASH", "BANG", "NEWLINE", "WS"
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
	public String getGrammarFileName() { return "HE.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public HEParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ProgramContext extends ParserRuleContext {
		public TerminalNode EOF() { return getToken(HEParser.EOF, 0); }
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
			if ( listener instanceof HEListener ) ((HEListener)listener).enterProgram(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitProgram(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitProgram(this);
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
			setState(97);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (((((_la - 11)) & ~0x3f) == 0 && ((1L << (_la - 11)) & 2307532272464664577L) != 0)) {
				{
				{
				setState(94);
				line();
				}
				}
				setState(99);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(100);
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
		public GlobalStatementContext globalStatement() {
			return getRuleContext(GlobalStatementContext.class,0);
		}
		public CommentContext comment() {
			return getRuleContext(CommentContext.class,0);
		}
		public TerminalNode NEWLINE() { return getToken(HEParser.NEWLINE, 0); }
		public LineContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_line; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterLine(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitLine(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitLine(this);
			else return visitor.visitChildren(this);
		}
	}

	public final LineContext line() throws RecognitionException {
		LineContext _localctx = new LineContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_line);
		try {
			setState(108);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,1,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(102);
				summon();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(103);
				withAssets();
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(104);
				object();
				}
				break;
			case 4:
				enterOuterAlt(_localctx, 4);
				{
				setState(105);
				globalStatement();
				}
				break;
			case 5:
				enterOuterAlt(_localctx, 5);
				{
				setState(106);
				comment();
				}
				break;
			case 6:
				enterOuterAlt(_localctx, 6);
				{
				setState(107);
				match(NEWLINE);
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
	public static class SummonContext extends ParserRuleContext {
		public TerminalNode SUMMON() { return getToken(HEParser.SUMMON, 0); }
		public TerminalNode STRING() { return getToken(HEParser.STRING, 0); }
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TerminalNode NAMED() { return getToken(HEParser.NAMED, 0); }
		public TerminalNode AS() { return getToken(HEParser.AS, 0); }
		public SummonContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_summon; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterSummon(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitSummon(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitSummon(this);
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
			setState(110);
			match(SUMMON);
			setState(111);
			match(STRING);
			setState(112);
			_la = _input.LA(1);
			if ( !(_la==NAMED || _la==AS) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(113);
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
		public TerminalNode WITH() { return getToken(HEParser.WITH, 0); }
		public AssetListContext assetList() {
			return getRuleContext(AssetListContext.class,0);
		}
		public WithAssetsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_withAssets; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterWithAssets(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitWithAssets(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitWithAssets(this);
			else return visitor.visitChildren(this);
		}
	}

	public final WithAssetsContext withAssets() throws RecognitionException {
		WithAssetsContext _localctx = new WithAssetsContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_withAssets);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(115);
			match(WITH);
			setState(116);
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
		public List<TerminalNode> AND() { return getTokens(HEParser.AND); }
		public TerminalNode AND(int i) {
			return getToken(HEParser.AND, i);
		}
		public AssetListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assetList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAssetList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAssetList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAssetList(this);
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
			setState(118);
			asset();
			setState(123);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==AND) {
				{
				{
				setState(119);
				match(AND);
				setState(120);
				asset();
				}
				}
				setState(125);
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
		public TerminalNode STRING() { return getToken(HEParser.STRING, 0); }
		public TerminalNode NAMED() { return getToken(HEParser.NAMED, 0); }
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public AssetContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_asset; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAsset(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAsset(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAsset(this);
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
			setState(126);
			assetType();
			setState(127);
			match(STRING);
			setState(130);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==NAMED) {
				{
				setState(128);
				match(NAMED);
				setState(129);
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
		public TerminalNode IMAGE() { return getToken(HEParser.IMAGE, 0); }
		public TerminalNode SOUND() { return getToken(HEParser.SOUND, 0); }
		public TerminalNode MUSIC() { return getToken(HEParser.MUSIC, 0); }
		public TerminalNode VIDEO() { return getToken(HEParser.VIDEO, 0); }
		public TerminalNode FONT() { return getToken(HEParser.FONT, 0); }
		public TerminalNode SHADER() { return getToken(HEParser.SHADER, 0); }
		public AssetTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assetType; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAssetType(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAssetType(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAssetType(this);
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
			setState(132);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 1134907106097364992L) != 0)) ) {
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
		public List<TerminalNode> ID() { return getTokens(HEParser.ID); }
		public TerminalNode ID(int i) {
			return getToken(HEParser.ID, i);
		}
		public ObjectBodyContext objectBody() {
			return getRuleContext(ObjectBodyContext.class,0);
		}
		public TerminalNode CREATE() { return getToken(HEParser.CREATE, 0); }
		public TerminalNode MAKE() { return getToken(HEParser.MAKE, 0); }
		public TerminalNode LIKE() { return getToken(HEParser.LIKE, 0); }
		public ObjectContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_object; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterObject(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitObject(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitObject(this);
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
			setState(134);
			_la = _input.LA(1);
			if ( !(_la==CREATE || _la==MAKE) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(135);
			match(ID);
			setState(138);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LIKE) {
				{
				setState(136);
				match(LIKE);
				setState(137);
				match(ID);
				}
			}

			setState(140);
			match(T__0);
			setState(141);
			objectBody();
			setState(142);
			match(T__1);
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
		public List<SectionContext> section() {
			return getRuleContexts(SectionContext.class);
		}
		public SectionContext section(int i) {
			return getRuleContext(SectionContext.class,i);
		}
		public List<ObjectContext> object() {
			return getRuleContexts(ObjectContext.class);
		}
		public ObjectContext object(int i) {
			return getRuleContext(ObjectContext.class,i);
		}
		public List<CommentContext> comment() {
			return getRuleContexts(CommentContext.class);
		}
		public CommentContext comment(int i) {
			return getRuleContext(CommentContext.class,i);
		}
		public List<WithAssetsContext> withAssets() {
			return getRuleContexts(WithAssetsContext.class);
		}
		public WithAssetsContext withAssets(int i) {
			return getRuleContext(WithAssetsContext.class,i);
		}
		public ObjectBodyContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_objectBody; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterObjectBody(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitObjectBody(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitObjectBody(this);
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
			setState(150);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 120469063680L) != 0)) {
				{
				setState(148);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case ID:
				case WHEN:
				case WHENEVER:
				case ON:
					{
					setState(144);
					section();
					}
					break;
				case CREATE:
				case MAKE:
					{
					setState(145);
					object();
					}
					break;
				case T__10:
					{
					setState(146);
					comment();
					}
					break;
				case WITH:
					{
					setState(147);
					withAssets();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				setState(152);
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
	public static class SectionContext extends ParserRuleContext {
		public PropertiesContext properties() {
			return getRuleContext(PropertiesContext.class,0);
		}
		public AbilitiesContext abilities() {
			return getRuleContext(AbilitiesContext.class,0);
		}
		public ReactionsContext reactions() {
			return getRuleContext(ReactionsContext.class,0);
		}
		public MemoriesContext memories() {
			return getRuleContext(MemoriesContext.class,0);
		}
		public SectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_section; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterSection(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitSection(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitSection(this);
			else return visitor.visitChildren(this);
		}
	}

	public final SectionContext section() throws RecognitionException {
		SectionContext _localctx = new SectionContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_section);
		try {
			setState(157);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,7,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(153);
				properties();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(154);
				abilities();
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(155);
				reactions();
				}
				break;
			case 4:
				enterOuterAlt(_localctx, 4);
				{
				setState(156);
				memories();
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
	public static class PropertiesContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public BlockContext block() {
			return getRuleContext(BlockContext.class,0);
		}
		public TerminalNode HAS() { return getToken(HEParser.HAS, 0); }
		public TerminalNode OWNS() { return getToken(HEParser.OWNS, 0); }
		public TerminalNode CARRIES() { return getToken(HEParser.CARRIES, 0); }
		public PropertiesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_properties; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterProperties(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitProperties(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitProperties(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PropertiesContext properties() throws RecognitionException {
		PropertiesContext _localctx = new PropertiesContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_properties);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(159);
			match(ID);
			setState(160);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 3758096384L) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(161);
			block();
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
	public static class AbilitiesContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TerminalNode CAN() { return getToken(HEParser.CAN, 0); }
		public TerminalNode KNOWS_HOW_TO() { return getToken(HEParser.KNOWS_HOW_TO, 0); }
		public List<ActionContext> action() {
			return getRuleContexts(ActionContext.class);
		}
		public ActionContext action(int i) {
			return getRuleContext(ActionContext.class,i);
		}
		public AbilitiesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_abilities; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAbilities(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAbilities(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAbilities(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AbilitiesContext abilities() throws RecognitionException {
		AbilitiesContext _localctx = new AbilitiesContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_abilities);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(163);
			match(ID);
			setState(164);
			_la = _input.LA(1);
			if ( !(_la==CAN || _la==KNOWS_HOW_TO) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(165);
			match(T__2);
			setState(169);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==ID) {
				{
				{
				setState(166);
				action();
				}
				}
				setState(171);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(172);
			match(T__3);
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
	public static class ReactionsContext extends ParserRuleContext {
		public TriggerContext trigger() {
			return getRuleContext(TriggerContext.class,0);
		}
		public TerminalNode WHEN() { return getToken(HEParser.WHEN, 0); }
		public TerminalNode WHENEVER() { return getToken(HEParser.WHENEVER, 0); }
		public TerminalNode ON() { return getToken(HEParser.ON, 0); }
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public ReactionsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_reactions; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterReactions(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitReactions(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitReactions(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ReactionsContext reactions() throws RecognitionException {
		ReactionsContext _localctx = new ReactionsContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_reactions);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(174);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 120259084288L) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(175);
			trigger();
			setState(176);
			match(T__2);
			setState(180);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3460956940140544000L) != 0)) {
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
			match(T__3);
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
	public static class MemoriesContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TerminalNode REMEMBERS() { return getToken(HEParser.REMEMBERS, 0); }
		public BlockContext block() {
			return getRuleContext(BlockContext.class,0);
		}
		public MemoriesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_memories; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterMemories(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitMemories(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitMemories(this);
			else return visitor.visitChildren(this);
		}
	}

	public final MemoriesContext memories() throws RecognitionException {
		MemoriesContext _localctx = new MemoriesContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_memories);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(185);
			match(ID);
			setState(186);
			match(REMEMBERS);
			setState(187);
			match(T__4);
			setState(188);
			block();
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
	public static class TriggerContext extends ParserRuleContext {
		public List<TerminalNode> ID() { return getTokens(HEParser.ID); }
		public TerminalNode ID(int i) {
			return getToken(HEParser.ID, i);
		}
		public TerminalNode WITH() { return getToken(HEParser.WITH, 0); }
		public TriggerContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_trigger; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterTrigger(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitTrigger(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitTrigger(this);
			else return visitor.visitChildren(this);
		}
	}

	public final TriggerContext trigger() throws RecognitionException {
		TriggerContext _localctx = new TriggerContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_trigger);
		try {
			setState(194);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,10,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(190);
				match(ID);
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(191);
				match(ID);
				setState(192);
				match(WITH);
				setState(193);
				match(ID);
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
	public static class BlockContext extends ParserRuleContext {
		public List<PropertyContext> property() {
			return getRuleContexts(PropertyContext.class);
		}
		public PropertyContext property(int i) {
			return getRuleContext(PropertyContext.class,i);
		}
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public List<CommentContext> comment() {
			return getRuleContexts(CommentContext.class);
		}
		public CommentContext comment(int i) {
			return getRuleContext(CommentContext.class,i);
		}
		public BlockContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_block; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterBlock(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitBlock(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitBlock(this);
			else return visitor.visitChildren(this);
		}
	}

	public final BlockContext block() throws RecognitionException {
		BlockContext _localctx = new BlockContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_block);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(196);
			match(T__2);
			setState(202);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3460956940140808192L) != 0)) {
				{
				setState(200);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case ID:
					{
					setState(197);
					property();
					}
					break;
				case MAKE:
				case SET:
				case IF:
				case REPEAT:
				case WHILE:
				case TELL:
				case WAIT:
				case RETURN:
				case SAY:
				case PRINT:
					{
					setState(198);
					statement();
					}
					break;
				case T__10:
					{
					setState(199);
					comment();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				setState(204);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(205);
			match(T__3);
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
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode IS() { return getToken(HEParser.IS, 0); }
		public TerminalNode STARTS_AS() { return getToken(HEParser.STARTS_AS, 0); }
		public PropertyContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_property; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterProperty(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitProperty(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitProperty(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PropertyContext property() throws RecognitionException {
		PropertyContext _localctx = new PropertyContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_property);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(207);
			match(ID);
			setState(208);
			_la = _input.LA(1);
			if ( !(_la==IS || _la==STARTS_AS) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(209);
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
	public static class ActionContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TerminalNode RETURNS() { return getToken(HEParser.RETURNS, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public ParamListContext paramList() {
			return getRuleContext(ParamListContext.class,0);
		}
		public ActionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_action; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAction(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAction(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAction(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ActionContext action() throws RecognitionException {
		ActionContext _localctx = new ActionContext(_ctx, getState());
		enterRule(_localctx, 34, RULE_action);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(211);
			match(ID);
			setState(217);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__5) {
				{
				setState(212);
				match(T__5);
				setState(214);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==ID) {
					{
					setState(213);
					paramList();
					}
				}

				setState(216);
				match(T__6);
				}
			}

			setState(221);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==RETURNS) {
				{
				setState(219);
				match(RETURNS);
				setState(220);
				type();
				}
			}

			setState(223);
			match(T__2);
			setState(227);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3460956940140544000L) != 0)) {
				{
				{
				setState(224);
				statement();
				}
				}
				setState(229);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(230);
			match(T__3);
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
		public ParamListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_paramList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterParamList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitParamList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitParamList(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ParamListContext paramList() throws RecognitionException {
		ParamListContext _localctx = new ParamListContext(_ctx, getState());
		enterRule(_localctx, 36, RULE_paramList);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(232);
			parameter();
			setState(237);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==T__7) {
				{
				{
				setState(233);
				match(T__7);
				setState(234);
				parameter();
				}
				}
				setState(239);
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
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public ParameterContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_parameter; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterParameter(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitParameter(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitParameter(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ParameterContext parameter() throws RecognitionException {
		ParameterContext _localctx = new ParameterContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_parameter);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(240);
			match(ID);
			setState(241);
			match(T__4);
			setState(242);
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
	public static class StatementContext extends ParserRuleContext {
		public SayContext say() {
			return getRuleContext(SayContext.class,0);
		}
		public ChangeContext change() {
			return getRuleContext(ChangeContext.class,0);
		}
		public DecideContext decide() {
			return getRuleContext(DecideContext.class,0);
		}
		public RepeatContext repeat() {
			return getRuleContext(RepeatContext.class,0);
		}
		public CallContext call() {
			return getRuleContext(CallContext.class,0);
		}
		public WaitStmtContext waitStmt() {
			return getRuleContext(WaitStmtContext.class,0);
		}
		public ReturnStmtContext returnStmt() {
			return getRuleContext(ReturnStmtContext.class,0);
		}
		public StatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_statement; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterStatement(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitStatement(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitStatement(this);
			else return visitor.visitChildren(this);
		}
	}

	public final StatementContext statement() throws RecognitionException {
		StatementContext _localctx = new StatementContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_statement);
		try {
			setState(251);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case SAY:
			case PRINT:
				enterOuterAlt(_localctx, 1);
				{
				setState(244);
				say();
				}
				break;
			case MAKE:
			case SET:
				enterOuterAlt(_localctx, 2);
				{
				setState(245);
				change();
				}
				break;
			case IF:
				enterOuterAlt(_localctx, 3);
				{
				setState(246);
				decide();
				}
				break;
			case REPEAT:
			case WHILE:
				enterOuterAlt(_localctx, 4);
				{
				setState(247);
				repeat();
				}
				break;
			case TELL:
				enterOuterAlt(_localctx, 5);
				{
				setState(248);
				call();
				}
				break;
			case WAIT:
				enterOuterAlt(_localctx, 6);
				{
				setState(249);
				waitStmt();
				}
				break;
			case RETURN:
				enterOuterAlt(_localctx, 7);
				{
				setState(250);
				returnStmt();
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
	public static class SayContext extends ParserRuleContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SAY() { return getToken(HEParser.SAY, 0); }
		public TerminalNode PRINT() { return getToken(HEParser.PRINT, 0); }
		public SayContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_say; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterSay(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitSay(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitSay(this);
			else return visitor.visitChildren(this);
		}
	}

	public final SayContext say() throws RecognitionException {
		SayContext _localctx = new SayContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_say);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(253);
			_la = _input.LA(1);
			if ( !(_la==SAY || _la==PRINT) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(254);
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
	public static class ChangeContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public TerminalNode TO() { return getToken(HEParser.TO, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SET() { return getToken(HEParser.SET, 0); }
		public TerminalNode MAKE() { return getToken(HEParser.MAKE, 0); }
		public ChangeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_change; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterChange(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitChange(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitChange(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ChangeContext change() throws RecognitionException {
		ChangeContext _localctx = new ChangeContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_change);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(256);
			_la = _input.LA(1);
			if ( !(_la==MAKE || _la==SET) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(257);
			match(ID);
			setState(258);
			match(TO);
			setState(259);
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
	public static class DecideContext extends ParserRuleContext {
		public TerminalNode IF() { return getToken(HEParser.IF, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode THEN() { return getToken(HEParser.THEN, 0); }
		public List<BlockContext> block() {
			return getRuleContexts(BlockContext.class);
		}
		public BlockContext block(int i) {
			return getRuleContext(BlockContext.class,i);
		}
		public TerminalNode ELSE() { return getToken(HEParser.ELSE, 0); }
		public DecideContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_decide; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterDecide(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitDecide(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitDecide(this);
			else return visitor.visitChildren(this);
		}
	}

	public final DecideContext decide() throws RecognitionException {
		DecideContext _localctx = new DecideContext(_ctx, getState());
		enterRule(_localctx, 46, RULE_decide);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(261);
			match(IF);
			setState(262);
			expression();
			setState(263);
			match(THEN);
			setState(264);
			block();
			setState(267);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ELSE) {
				{
				setState(265);
				match(ELSE);
				setState(266);
				block();
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
	public static class RepeatContext extends ParserRuleContext {
		public BlockContext block() {
			return getRuleContext(BlockContext.class,0);
		}
		public TerminalNode REPEAT() { return getToken(HEParser.REPEAT, 0); }
		public TerminalNode WHILE() { return getToken(HEParser.WHILE, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public RepeatContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_repeat; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterRepeat(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitRepeat(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitRepeat(this);
			else return visitor.visitChildren(this);
		}
	}

	public final RepeatContext repeat() throws RecognitionException {
		RepeatContext _localctx = new RepeatContext(_ctx, getState());
		enterRule(_localctx, 48, RULE_repeat);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(272);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case REPEAT:
				{
				setState(269);
				match(REPEAT);
				}
				break;
			case WHILE:
				{
				setState(270);
				match(WHILE);
				setState(271);
				expression();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(274);
			block();
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
	public static class CallContext extends ParserRuleContext {
		public TerminalNode TELL() { return getToken(HEParser.TELL, 0); }
		public List<TerminalNode> ID() { return getTokens(HEParser.ID); }
		public TerminalNode ID(int i) {
			return getToken(HEParser.ID, i);
		}
		public TerminalNode TO() { return getToken(HEParser.TO, 0); }
		public TerminalNode WITH() { return getToken(HEParser.WITH, 0); }
		public ArgListContext argList() {
			return getRuleContext(ArgListContext.class,0);
		}
		public CallContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_call; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterCall(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitCall(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitCall(this);
			else return visitor.visitChildren(this);
		}
	}

	public final CallContext call() throws RecognitionException {
		CallContext _localctx = new CallContext(_ctx, getState());
		enterRule(_localctx, 50, RULE_call);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(276);
			match(TELL);
			setState(277);
			match(ID);
			setState(278);
			match(TO);
			setState(279);
			match(ID);
			setState(282);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,21,_ctx) ) {
			case 1:
				{
				setState(280);
				match(WITH);
				setState(281);
				argList();
				}
				break;
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
		public TerminalNode WAIT() { return getToken(HEParser.WAIT, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SECONDS() { return getToken(HEParser.SECONDS, 0); }
		public TerminalNode FRAMES() { return getToken(HEParser.FRAMES, 0); }
		public WaitStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_waitStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterWaitStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitWaitStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitWaitStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final WaitStmtContext waitStmt() throws RecognitionException {
		WaitStmtContext _localctx = new WaitStmtContext(_ctx, getState());
		enterRule(_localctx, 52, RULE_waitStmt);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(284);
			match(WAIT);
			setState(285);
			expression();
			setState(287);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==SECONDS || _la==FRAMES) {
				{
				setState(286);
				_la = _input.LA(1);
				if ( !(_la==SECONDS || _la==FRAMES) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
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
	public static class ReturnStmtContext extends ParserRuleContext {
		public TerminalNode RETURN() { return getToken(HEParser.RETURN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ReturnStmtContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_returnStmt; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterReturnStmt(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitReturnStmt(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitReturnStmt(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ReturnStmtContext returnStmt() throws RecognitionException {
		ReturnStmtContext _localctx = new ReturnStmtContext(_ctx, getState());
		enterRule(_localctx, 54, RULE_returnStmt);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(289);
			match(RETURN);
			setState(290);
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
		public ArgListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_argList; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterArgList(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitArgList(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitArgList(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ArgListContext argList() throws RecognitionException {
		ArgListContext _localctx = new ArgListContext(_ctx, getState());
		enterRule(_localctx, 56, RULE_argList);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(292);
			expression();
			setState(297);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==T__7) {
				{
				{
				setState(293);
				match(T__7);
				setState(294);
				expression();
				}
				}
				setState(299);
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
		public LogicOrContext logicOr() {
			return getRuleContext(LogicOrContext.class,0);
		}
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ExpressionContext expression() throws RecognitionException {
		ExpressionContext _localctx = new ExpressionContext(_ctx, getState());
		enterRule(_localctx, 58, RULE_expression);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(300);
			logicOr();
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
	public static class LogicOrContext extends ParserRuleContext {
		public List<LogicAndContext> logicAnd() {
			return getRuleContexts(LogicAndContext.class);
		}
		public LogicAndContext logicAnd(int i) {
			return getRuleContext(LogicAndContext.class,i);
		}
		public List<TerminalNode> OR() { return getTokens(HEParser.OR); }
		public TerminalNode OR(int i) {
			return getToken(HEParser.OR, i);
		}
		public LogicOrContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_logicOr; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterLogicOr(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitLogicOr(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitLogicOr(this);
			else return visitor.visitChildren(this);
		}
	}

	public final LogicOrContext logicOr() throws RecognitionException {
		LogicOrContext _localctx = new LogicOrContext(_ctx, getState());
		enterRule(_localctx, 60, RULE_logicOr);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(302);
			logicAnd();
			setState(307);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==OR) {
				{
				{
				setState(303);
				match(OR);
				setState(304);
				logicAnd();
				}
				}
				setState(309);
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
	public static class LogicAndContext extends ParserRuleContext {
		public List<ComparisonContext> comparison() {
			return getRuleContexts(ComparisonContext.class);
		}
		public ComparisonContext comparison(int i) {
			return getRuleContext(ComparisonContext.class,i);
		}
		public List<TerminalNode> AND() { return getTokens(HEParser.AND); }
		public TerminalNode AND(int i) {
			return getToken(HEParser.AND, i);
		}
		public LogicAndContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_logicAnd; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterLogicAnd(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitLogicAnd(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitLogicAnd(this);
			else return visitor.visitChildren(this);
		}
	}

	public final LogicAndContext logicAnd() throws RecognitionException {
		LogicAndContext _localctx = new LogicAndContext(_ctx, getState());
		enterRule(_localctx, 62, RULE_logicAnd);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(310);
			comparison();
			setState(315);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==AND) {
				{
				{
				setState(311);
				match(AND);
				setState(312);
				comparison();
				}
				}
				setState(317);
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
	public static class ComparisonContext extends ParserRuleContext {
		public List<ArithmeticContext> arithmetic() {
			return getRuleContexts(ArithmeticContext.class);
		}
		public ArithmeticContext arithmetic(int i) {
			return getRuleContext(ArithmeticContext.class,i);
		}
		public CompOpContext compOp() {
			return getRuleContext(CompOpContext.class,0);
		}
		public ComparisonContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_comparison; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterComparison(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitComparison(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitComparison(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ComparisonContext comparison() throws RecognitionException {
		ComparisonContext _localctx = new ComparisonContext(_ctx, getState());
		enterRule(_localctx, 64, RULE_comparison);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(318);
			arithmetic();
			setState(322);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 258048L) != 0)) {
				{
				setState(319);
				compOp();
				setState(320);
				arithmetic();
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
	public static class ArithmeticContext extends ParserRuleContext {
		public List<TermContext> term() {
			return getRuleContexts(TermContext.class);
		}
		public TermContext term(int i) {
			return getRuleContext(TermContext.class,i);
		}
		public List<AddOpContext> addOp() {
			return getRuleContexts(AddOpContext.class);
		}
		public AddOpContext addOp(int i) {
			return getRuleContext(AddOpContext.class,i);
		}
		public ArithmeticContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arithmetic; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterArithmetic(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitArithmetic(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitArithmetic(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ArithmeticContext arithmetic() throws RecognitionException {
		ArithmeticContext _localctx = new ArithmeticContext(_ctx, getState());
		enterRule(_localctx, 66, RULE_arithmetic);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(324);
			term();
			setState(330);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==PLUS || _la==MINUS) {
				{
				{
				setState(325);
				addOp();
				setState(326);
				term();
				}
				}
				setState(332);
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
	public static class TermContext extends ParserRuleContext {
		public List<FactorContext> factor() {
			return getRuleContexts(FactorContext.class);
		}
		public FactorContext factor(int i) {
			return getRuleContext(FactorContext.class,i);
		}
		public List<MultOpContext> multOp() {
			return getRuleContexts(MultOpContext.class);
		}
		public MultOpContext multOp(int i) {
			return getRuleContext(MultOpContext.class,i);
		}
		public TermContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_term; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterTerm(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitTerm(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitTerm(this);
			else return visitor.visitChildren(this);
		}
	}

	public final TermContext term() throws RecognitionException {
		TermContext _localctx = new TermContext(_ctx, getState());
		enterRule(_localctx, 68, RULE_term);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(333);
			factor();
			setState(339);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==STAR || _la==SLASH) {
				{
				{
				setState(334);
				multOp();
				setState(335);
				factor();
				}
				}
				setState(341);
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
	public static class FactorContext extends ParserRuleContext {
		public UnaryContext unary() {
			return getRuleContext(UnaryContext.class,0);
		}
		public FactorContext factor() {
			return getRuleContext(FactorContext.class,0);
		}
		public FactorContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_factor; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterFactor(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitFactor(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitFactor(this);
			else return visitor.visitChildren(this);
		}
	}

	public final FactorContext factor() throws RecognitionException {
		FactorContext _localctx = new FactorContext(_ctx, getState());
		enterRule(_localctx, 70, RULE_factor);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(342);
			unary();
			setState(345);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__8) {
				{
				setState(343);
				match(T__8);
				setState(344);
				factor();
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
	public static class UnaryContext extends ParserRuleContext {
		public PrimaryContext primary() {
			return getRuleContext(PrimaryContext.class,0);
		}
		public TerminalNode MINUS() { return getToken(HEParser.MINUS, 0); }
		public TerminalNode BANG() { return getToken(HEParser.BANG, 0); }
		public UnaryContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_unary; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterUnary(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitUnary(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitUnary(this);
			else return visitor.visitChildren(this);
		}
	}

	public final UnaryContext unary() throws RecognitionException {
		UnaryContext _localctx = new UnaryContext(_ctx, getState());
		enterRule(_localctx, 72, RULE_unary);
		int _la;
		try {
			setState(350);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case MINUS:
			case BANG:
				enterOuterAlt(_localctx, 1);
				{
				setState(347);
				_la = _input.LA(1);
				if ( !(_la==MINUS || _la==BANG) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(348);
				primary();
				}
				break;
			case T__2:
			case T__5:
			case ID:
			case NUMBER:
			case STRING:
			case BOOLEAN:
				enterOuterAlt(_localctx, 2);
				{
				setState(349);
				primary();
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
	public static class PrimaryContext extends ParserRuleContext {
		public TerminalNode NUMBER() { return getToken(HEParser.NUMBER, 0); }
		public TerminalNode STRING() { return getToken(HEParser.STRING, 0); }
		public TerminalNode BOOLEAN() { return getToken(HEParser.BOOLEAN, 0); }
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public CallExpressionContext callExpression() {
			return getRuleContext(CallExpressionContext.class,0);
		}
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ArrayLiteralContext arrayLiteral() {
			return getRuleContext(ArrayLiteralContext.class,0);
		}
		public PrimaryContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_primary; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterPrimary(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitPrimary(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitPrimary(this);
			else return visitor.visitChildren(this);
		}
	}

	public final PrimaryContext primary() throws RecognitionException {
		PrimaryContext _localctx = new PrimaryContext(_ctx, getState());
		enterRule(_localctx, 74, RULE_primary);
		try {
			setState(362);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,31,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(352);
				match(NUMBER);
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(353);
				match(STRING);
				}
				break;
			case 3:
				enterOuterAlt(_localctx, 3);
				{
				setState(354);
				match(BOOLEAN);
				}
				break;
			case 4:
				enterOuterAlt(_localctx, 4);
				{
				setState(355);
				match(ID);
				}
				break;
			case 5:
				enterOuterAlt(_localctx, 5);
				{
				setState(356);
				callExpression();
				}
				break;
			case 6:
				enterOuterAlt(_localctx, 6);
				{
				setState(357);
				match(T__5);
				setState(358);
				expression();
				setState(359);
				match(T__6);
				}
				break;
			case 7:
				enterOuterAlt(_localctx, 7);
				{
				setState(361);
				arrayLiteral();
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
	public static class CallExpressionContext extends ParserRuleContext {
		public TerminalNode ID() { return getToken(HEParser.ID, 0); }
		public ArgListContext argList() {
			return getRuleContext(ArgListContext.class,0);
		}
		public CallExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_callExpression; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterCallExpression(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitCallExpression(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitCallExpression(this);
			else return visitor.visitChildren(this);
		}
	}

	public final CallExpressionContext callExpression() throws RecognitionException {
		CallExpressionContext _localctx = new CallExpressionContext(_ctx, getState());
		enterRule(_localctx, 76, RULE_callExpression);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(364);
			match(ID);
			setState(365);
			match(T__5);
			setState(367);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3932232L) != 0) || _la==MINUS || _la==BANG) {
				{
				setState(366);
				argList();
				}
			}

			setState(369);
			match(T__6);
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
	public static class ArrayLiteralContext extends ParserRuleContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public ArrayLiteralContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arrayLiteral; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterArrayLiteral(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitArrayLiteral(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitArrayLiteral(this);
			else return visitor.visitChildren(this);
		}
	}

	public final ArrayLiteralContext arrayLiteral() throws RecognitionException {
		ArrayLiteralContext _localctx = new ArrayLiteralContext(_ctx, getState());
		enterRule(_localctx, 78, RULE_arrayLiteral);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(371);
			match(T__2);
			setState(380);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3932232L) != 0) || _la==MINUS || _la==BANG) {
				{
				setState(372);
				expression();
				setState(377);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==T__7) {
					{
					{
					setState(373);
					match(T__7);
					setState(374);
					expression();
					}
					}
					setState(379);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				}
			}

			setState(382);
			match(T__3);
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
		public BaseTypeContext baseType() {
			return getRuleContext(BaseTypeContext.class,0);
		}
		public TypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_type; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterType(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitType(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitType(this);
			else return visitor.visitChildren(this);
		}
	}

	public final TypeContext type() throws RecognitionException {
		TypeContext _localctx = new TypeContext(_ctx, getState());
		enterRule(_localctx, 80, RULE_type);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(384);
			baseType();
			setState(386);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__9) {
				{
				setState(385);
				match(T__9);
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
	public static class BaseTypeContext extends ParserRuleContext {
		public TerminalNode NUMBER_TYPE() { return getToken(HEParser.NUMBER_TYPE, 0); }
		public TerminalNode STRING_TYPE() { return getToken(HEParser.STRING_TYPE, 0); }
		public TerminalNode BOOLEAN_TYPE() { return getToken(HEParser.BOOLEAN_TYPE, 0); }
		public BaseTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_baseType; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterBaseType(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitBaseType(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitBaseType(this);
			else return visitor.visitChildren(this);
		}
	}

	public final BaseTypeContext baseType() throws RecognitionException {
		BaseTypeContext _localctx = new BaseTypeContext(_ctx, getState());
		enterRule(_localctx, 82, RULE_baseType);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(388);
			_la = _input.LA(1);
			if ( !(((((_la - 64)) & ~0x3f) == 0 && ((1L << (_la - 64)) & 7L) != 0)) ) {
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
	public static class GlobalStatementContext extends ParserRuleContext {
		public SayContext say() {
			return getRuleContext(SayContext.class,0);
		}
		public ChangeContext change() {
			return getRuleContext(ChangeContext.class,0);
		}
		public WaitStmtContext waitStmt() {
			return getRuleContext(WaitStmtContext.class,0);
		}
		public CallContext call() {
			return getRuleContext(CallContext.class,0);
		}
		public GlobalStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_globalStatement; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterGlobalStatement(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitGlobalStatement(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitGlobalStatement(this);
			else return visitor.visitChildren(this);
		}
	}

	public final GlobalStatementContext globalStatement() throws RecognitionException {
		GlobalStatementContext _localctx = new GlobalStatementContext(_ctx, getState());
		enterRule(_localctx, 84, RULE_globalStatement);
		try {
			setState(394);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case SAY:
			case PRINT:
				enterOuterAlt(_localctx, 1);
				{
				setState(390);
				say();
				}
				break;
			case MAKE:
			case SET:
				enterOuterAlt(_localctx, 2);
				{
				setState(391);
				change();
				}
				break;
			case WAIT:
				enterOuterAlt(_localctx, 3);
				{
				setState(392);
				waitStmt();
				}
				break;
			case TELL:
				enterOuterAlt(_localctx, 4);
				{
				setState(393);
				call();
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
	public static class CommentContext extends ParserRuleContext {
		public CommentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_comment; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterComment(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitComment(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitComment(this);
			else return visitor.visitChildren(this);
		}
	}

	public final CommentContext comment() throws RecognitionException {
		CommentContext _localctx = new CommentContext(_ctx, getState());
		enterRule(_localctx, 86, RULE_comment);
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(396);
			match(T__10);
			setState(400);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,37,_ctx);
			while ( _alt!=1 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1+1 ) {
					{
					{
					setState(397);
					matchWildcard();
					}
					} 
				}
				setState(402);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,37,_ctx);
			}
			setState(403);
			match(T__10);
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
	public static class CompOpContext extends ParserRuleContext {
		public CompOpContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_compOp; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterCompOp(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitCompOp(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitCompOp(this);
			else return visitor.visitChildren(this);
		}
	}

	public final CompOpContext compOp() throws RecognitionException {
		CompOpContext _localctx = new CompOpContext(_ctx, getState());
		enterRule(_localctx, 88, RULE_compOp);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(405);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & 258048L) != 0)) ) {
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
	public static class AddOpContext extends ParserRuleContext {
		public TerminalNode PLUS() { return getToken(HEParser.PLUS, 0); }
		public TerminalNode MINUS() { return getToken(HEParser.MINUS, 0); }
		public AddOpContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_addOp; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterAddOp(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitAddOp(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitAddOp(this);
			else return visitor.visitChildren(this);
		}
	}

	public final AddOpContext addOp() throws RecognitionException {
		AddOpContext _localctx = new AddOpContext(_ctx, getState());
		enterRule(_localctx, 90, RULE_addOp);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(407);
			_la = _input.LA(1);
			if ( !(_la==PLUS || _la==MINUS) ) {
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
	public static class MultOpContext extends ParserRuleContext {
		public TerminalNode STAR() { return getToken(HEParser.STAR, 0); }
		public TerminalNode SLASH() { return getToken(HEParser.SLASH, 0); }
		public MultOpContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_multOp; }
		@Override
		public void enterRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).enterMultOp(this);
		}
		@Override
		public void exitRule(ParseTreeListener listener) {
			if ( listener instanceof HEListener ) ((HEListener)listener).exitMultOp(this);
		}
		@Override
		public <T> T accept(ParseTreeVisitor<? extends T> visitor) {
			if ( visitor instanceof HEVisitor ) return ((HEVisitor<? extends T>)visitor).visitMultOp(this);
			else return visitor.visitChildren(this);
		}
	}

	public final MultOpContext multOp() throws RecognitionException {
		MultOpContext _localctx = new MultOpContext(_ctx, getState());
		enterRule(_localctx, 92, RULE_multOp);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(409);
			_la = _input.LA(1);
			if ( !(_la==STAR || _la==SLASH) ) {
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

	public static final String _serializedATN =
		"\u0004\u0001I\u019c\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001\u0002"+
		"\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004\u0002"+
		"\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007\u0002"+
		"\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b\u0002"+
		"\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007\u000f"+
		"\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0002\u0012\u0007\u0012"+
		"\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014\u0002\u0015\u0007\u0015"+
		"\u0002\u0016\u0007\u0016\u0002\u0017\u0007\u0017\u0002\u0018\u0007\u0018"+
		"\u0002\u0019\u0007\u0019\u0002\u001a\u0007\u001a\u0002\u001b\u0007\u001b"+
		"\u0002\u001c\u0007\u001c\u0002\u001d\u0007\u001d\u0002\u001e\u0007\u001e"+
		"\u0002\u001f\u0007\u001f\u0002 \u0007 \u0002!\u0007!\u0002\"\u0007\"\u0002"+
		"#\u0007#\u0002$\u0007$\u0002%\u0007%\u0002&\u0007&\u0002\'\u0007\'\u0002"+
		"(\u0007(\u0002)\u0007)\u0002*\u0007*\u0002+\u0007+\u0002,\u0007,\u0002"+
		"-\u0007-\u0002.\u0007.\u0001\u0000\u0005\u0000`\b\u0000\n\u0000\f\u0000"+
		"c\t\u0000\u0001\u0000\u0001\u0000\u0001\u0001\u0001\u0001\u0001\u0001"+
		"\u0001\u0001\u0001\u0001\u0001\u0001\u0003\u0001m\b\u0001\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0003\u0001\u0003"+
		"\u0001\u0003\u0001\u0004\u0001\u0004\u0001\u0004\u0005\u0004z\b\u0004"+
		"\n\u0004\f\u0004}\t\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005"+
		"\u0003\u0005\u0083\b\u0005\u0001\u0006\u0001\u0006\u0001\u0007\u0001\u0007"+
		"\u0001\u0007\u0001\u0007\u0003\u0007\u008b\b\u0007\u0001\u0007\u0001\u0007"+
		"\u0001\u0007\u0001\u0007\u0001\b\u0001\b\u0001\b\u0001\b\u0005\b\u0095"+
		"\b\b\n\b\f\b\u0098\t\b\u0001\t\u0001\t\u0001\t\u0001\t\u0003\t\u009e\b"+
		"\t\u0001\n\u0001\n\u0001\n\u0001\n\u0001\u000b\u0001\u000b\u0001\u000b"+
		"\u0001\u000b\u0005\u000b\u00a8\b\u000b\n\u000b\f\u000b\u00ab\t\u000b\u0001"+
		"\u000b\u0001\u000b\u0001\f\u0001\f\u0001\f\u0001\f\u0005\f\u00b3\b\f\n"+
		"\f\f\f\u00b6\t\f\u0001\f\u0001\f\u0001\r\u0001\r\u0001\r\u0001\r\u0001"+
		"\r\u0001\u000e\u0001\u000e\u0001\u000e\u0001\u000e\u0003\u000e\u00c3\b"+
		"\u000e\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0005\u000f\u00c9"+
		"\b\u000f\n\u000f\f\u000f\u00cc\t\u000f\u0001\u000f\u0001\u000f\u0001\u0010"+
		"\u0001\u0010\u0001\u0010\u0001\u0010\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0003\u0011\u00d7\b\u0011\u0001\u0011\u0003\u0011\u00da\b\u0011\u0001"+
		"\u0011\u0001\u0011\u0003\u0011\u00de\b\u0011\u0001\u0011\u0001\u0011\u0005"+
		"\u0011\u00e2\b\u0011\n\u0011\f\u0011\u00e5\t\u0011\u0001\u0011\u0001\u0011"+
		"\u0001\u0012\u0001\u0012\u0001\u0012\u0005\u0012\u00ec\b\u0012\n\u0012"+
		"\f\u0012\u00ef\t\u0012\u0001\u0013\u0001\u0013\u0001\u0013\u0001\u0013"+
		"\u0001\u0014\u0001\u0014\u0001\u0014\u0001\u0014\u0001\u0014\u0001\u0014"+
		"\u0001\u0014\u0003\u0014\u00fc\b\u0014\u0001\u0015\u0001\u0015\u0001\u0015"+
		"\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0017"+
		"\u0001\u0017\u0001\u0017\u0001\u0017\u0001\u0017\u0001\u0017\u0003\u0017"+
		"\u010c\b\u0017\u0001\u0018\u0001\u0018\u0001\u0018\u0003\u0018\u0111\b"+
		"\u0018\u0001\u0018\u0001\u0018\u0001\u0019\u0001\u0019\u0001\u0019\u0001"+
		"\u0019\u0001\u0019\u0001\u0019\u0003\u0019\u011b\b\u0019\u0001\u001a\u0001"+
		"\u001a\u0001\u001a\u0003\u001a\u0120\b\u001a\u0001\u001b\u0001\u001b\u0001"+
		"\u001b\u0001\u001c\u0001\u001c\u0001\u001c\u0005\u001c\u0128\b\u001c\n"+
		"\u001c\f\u001c\u012b\t\u001c\u0001\u001d\u0001\u001d\u0001\u001e\u0001"+
		"\u001e\u0001\u001e\u0005\u001e\u0132\b\u001e\n\u001e\f\u001e\u0135\t\u001e"+
		"\u0001\u001f\u0001\u001f\u0001\u001f\u0005\u001f\u013a\b\u001f\n\u001f"+
		"\f\u001f\u013d\t\u001f\u0001 \u0001 \u0001 \u0001 \u0003 \u0143\b \u0001"+
		"!\u0001!\u0001!\u0001!\u0005!\u0149\b!\n!\f!\u014c\t!\u0001\"\u0001\""+
		"\u0001\"\u0001\"\u0005\"\u0152\b\"\n\"\f\"\u0155\t\"\u0001#\u0001#\u0001"+
		"#\u0003#\u015a\b#\u0001$\u0001$\u0001$\u0003$\u015f\b$\u0001%\u0001%\u0001"+
		"%\u0001%\u0001%\u0001%\u0001%\u0001%\u0001%\u0001%\u0003%\u016b\b%\u0001"+
		"&\u0001&\u0001&\u0003&\u0170\b&\u0001&\u0001&\u0001\'\u0001\'\u0001\'"+
		"\u0001\'\u0005\'\u0178\b\'\n\'\f\'\u017b\t\'\u0003\'\u017d\b\'\u0001\'"+
		"\u0001\'\u0001(\u0001(\u0003(\u0183\b(\u0001)\u0001)\u0001*\u0001*\u0001"+
		"*\u0001*\u0003*\u018b\b*\u0001+\u0001+\u0005+\u018f\b+\n+\f+\u0192\t+"+
		"\u0001+\u0001+\u0001,\u0001,\u0001-\u0001-\u0001.\u0001.\u0001.\u0001"+
		"\u0190\u0000/\u0000\u0002\u0004\u0006\b\n\f\u000e\u0010\u0012\u0014\u0016"+
		"\u0018\u001a\u001c\u001e \"$&(*,.02468:<>@BDFHJLNPRTVXZ\\\u0000\u000f"+
		"\u0001\u0000\u0018\u0019\u0001\u00006;\u0001\u0000\u001a\u001b\u0001\u0000"+
		"\u001d\u001f\u0001\u0000 !\u0001\u0000\"$\u0001\u0000&\'\u0001\u0000<"+
		"=\u0002\u0000\u001b\u001b))\u0001\u0000>?\u0002\u0000DDGG\u0001\u0000"+
		"@B\u0001\u0000\f\u0011\u0001\u0000CD\u0001\u0000EF\u01a7\u0000a\u0001"+
		"\u0000\u0000\u0000\u0002l\u0001\u0000\u0000\u0000\u0004n\u0001\u0000\u0000"+
		"\u0000\u0006s\u0001\u0000\u0000\u0000\bv\u0001\u0000\u0000\u0000\n~\u0001"+
		"\u0000\u0000\u0000\f\u0084\u0001\u0000\u0000\u0000\u000e\u0086\u0001\u0000"+
		"\u0000\u0000\u0010\u0096\u0001\u0000\u0000\u0000\u0012\u009d\u0001\u0000"+
		"\u0000\u0000\u0014\u009f\u0001\u0000\u0000\u0000\u0016\u00a3\u0001\u0000"+
		"\u0000\u0000\u0018\u00ae\u0001\u0000\u0000\u0000\u001a\u00b9\u0001\u0000"+
		"\u0000\u0000\u001c\u00c2\u0001\u0000\u0000\u0000\u001e\u00c4\u0001\u0000"+
		"\u0000\u0000 \u00cf\u0001\u0000\u0000\u0000\"\u00d3\u0001\u0000\u0000"+
		"\u0000$\u00e8\u0001\u0000\u0000\u0000&\u00f0\u0001\u0000\u0000\u0000("+
		"\u00fb\u0001\u0000\u0000\u0000*\u00fd\u0001\u0000\u0000\u0000,\u0100\u0001"+
		"\u0000\u0000\u0000.\u0105\u0001\u0000\u0000\u00000\u0110\u0001\u0000\u0000"+
		"\u00002\u0114\u0001\u0000\u0000\u00004\u011c\u0001\u0000\u0000\u00006"+
		"\u0121\u0001\u0000\u0000\u00008\u0124\u0001\u0000\u0000\u0000:\u012c\u0001"+
		"\u0000\u0000\u0000<\u012e\u0001\u0000\u0000\u0000>\u0136\u0001\u0000\u0000"+
		"\u0000@\u013e\u0001\u0000\u0000\u0000B\u0144\u0001\u0000\u0000\u0000D"+
		"\u014d\u0001\u0000\u0000\u0000F\u0156\u0001\u0000\u0000\u0000H\u015e\u0001"+
		"\u0000\u0000\u0000J\u016a\u0001\u0000\u0000\u0000L\u016c\u0001\u0000\u0000"+
		"\u0000N\u0173\u0001\u0000\u0000\u0000P\u0180\u0001\u0000\u0000\u0000R"+
		"\u0184\u0001\u0000\u0000\u0000T\u018a\u0001\u0000\u0000\u0000V\u018c\u0001"+
		"\u0000\u0000\u0000X\u0195\u0001\u0000\u0000\u0000Z\u0197\u0001\u0000\u0000"+
		"\u0000\\\u0199\u0001\u0000\u0000\u0000^`\u0003\u0002\u0001\u0000_^\u0001"+
		"\u0000\u0000\u0000`c\u0001\u0000\u0000\u0000a_\u0001\u0000\u0000\u0000"+
		"ab\u0001\u0000\u0000\u0000bd\u0001\u0000\u0000\u0000ca\u0001\u0000\u0000"+
		"\u0000de\u0005\u0000\u0000\u0001e\u0001\u0001\u0000\u0000\u0000fm\u0003"+
		"\u0004\u0002\u0000gm\u0003\u0006\u0003\u0000hm\u0003\u000e\u0007\u0000"+
		"im\u0003T*\u0000jm\u0003V+\u0000km\u0005H\u0000\u0000lf\u0001\u0000\u0000"+
		"\u0000lg\u0001\u0000\u0000\u0000lh\u0001\u0000\u0000\u0000li\u0001\u0000"+
		"\u0000\u0000lj\u0001\u0000\u0000\u0000lk\u0001\u0000\u0000\u0000m\u0003"+
		"\u0001\u0000\u0000\u0000no\u0005\u0016\u0000\u0000op\u0005\u0014\u0000"+
		"\u0000pq\u0007\u0000\u0000\u0000qr\u0005\u0012\u0000\u0000r\u0005\u0001"+
		"\u0000\u0000\u0000st\u0005\u0017\u0000\u0000tu\u0003\b\u0004\u0000u\u0007"+
		"\u0001\u0000\u0000\u0000v{\u0003\n\u0005\u0000wx\u00054\u0000\u0000xz"+
		"\u0003\n\u0005\u0000yw\u0001\u0000\u0000\u0000z}\u0001\u0000\u0000\u0000"+
		"{y\u0001\u0000\u0000\u0000{|\u0001\u0000\u0000\u0000|\t\u0001\u0000\u0000"+
		"\u0000}{\u0001\u0000\u0000\u0000~\u007f\u0003\f\u0006\u0000\u007f\u0082"+
		"\u0005\u0014\u0000\u0000\u0080\u0081\u0005\u0018\u0000\u0000\u0081\u0083"+
		"\u0005\u0012\u0000\u0000\u0082\u0080\u0001\u0000\u0000\u0000\u0082\u0083"+
		"\u0001\u0000\u0000\u0000\u0083\u000b\u0001\u0000\u0000\u0000\u0084\u0085"+
		"\u0007\u0001\u0000\u0000\u0085\r\u0001\u0000\u0000\u0000\u0086\u0087\u0007"+
		"\u0002\u0000\u0000\u0087\u008a\u0005\u0012\u0000\u0000\u0088\u0089\u0005"+
		"\u001c\u0000\u0000\u0089\u008b\u0005\u0012\u0000\u0000\u008a\u0088\u0001"+
		"\u0000\u0000\u0000\u008a\u008b\u0001\u0000\u0000\u0000\u008b\u008c\u0001"+
		"\u0000\u0000\u0000\u008c\u008d\u0005\u0001\u0000\u0000\u008d\u008e\u0003"+
		"\u0010\b\u0000\u008e\u008f\u0005\u0002\u0000\u0000\u008f\u000f\u0001\u0000"+
		"\u0000\u0000\u0090\u0095\u0003\u0012\t\u0000\u0091\u0095\u0003\u000e\u0007"+
		"\u0000\u0092\u0095\u0003V+\u0000\u0093\u0095\u0003\u0006\u0003\u0000\u0094"+
		"\u0090\u0001\u0000\u0000\u0000\u0094\u0091\u0001\u0000\u0000\u0000\u0094"+
		"\u0092\u0001\u0000\u0000\u0000\u0094\u0093\u0001\u0000\u0000\u0000\u0095"+
		"\u0098\u0001\u0000\u0000\u0000\u0096\u0094\u0001\u0000\u0000\u0000\u0096"+
		"\u0097\u0001\u0000\u0000\u0000\u0097\u0011\u0001\u0000\u0000\u0000\u0098"+
		"\u0096\u0001\u0000\u0000\u0000\u0099\u009e\u0003\u0014\n\u0000\u009a\u009e"+
		"\u0003\u0016\u000b\u0000\u009b\u009e\u0003\u0018\f\u0000\u009c\u009e\u0003"+
		"\u001a\r\u0000\u009d\u0099\u0001\u0000\u0000\u0000\u009d\u009a\u0001\u0000"+
		"\u0000\u0000\u009d\u009b\u0001\u0000\u0000\u0000\u009d\u009c\u0001\u0000"+
		"\u0000\u0000\u009e\u0013\u0001\u0000\u0000\u0000\u009f\u00a0\u0005\u0012"+
		"\u0000\u0000\u00a0\u00a1\u0007\u0003\u0000\u0000\u00a1\u00a2\u0003\u001e"+
		"\u000f\u0000\u00a2\u0015\u0001\u0000\u0000\u0000\u00a3\u00a4\u0005\u0012"+
		"\u0000\u0000\u00a4\u00a5\u0007\u0004\u0000\u0000\u00a5\u00a9\u0005\u0003"+
		"\u0000\u0000\u00a6\u00a8\u0003\"\u0011\u0000\u00a7\u00a6\u0001\u0000\u0000"+
		"\u0000\u00a8\u00ab\u0001\u0000\u0000\u0000\u00a9\u00a7\u0001\u0000\u0000"+
		"\u0000\u00a9\u00aa\u0001\u0000\u0000\u0000\u00aa\u00ac\u0001\u0000\u0000"+
		"\u0000\u00ab\u00a9\u0001\u0000\u0000\u0000\u00ac\u00ad\u0005\u0004\u0000"+
		"\u0000\u00ad\u0017\u0001\u0000\u0000\u0000\u00ae\u00af\u0007\u0005\u0000"+
		"\u0000\u00af\u00b0\u0003\u001c\u000e\u0000\u00b0\u00b4\u0005\u0003\u0000"+
		"\u0000\u00b1\u00b3\u0003(\u0014\u0000\u00b2\u00b1\u0001\u0000\u0000\u0000"+
		"\u00b3\u00b6\u0001\u0000\u0000\u0000\u00b4\u00b2\u0001\u0000\u0000\u0000"+
		"\u00b4\u00b5\u0001\u0000\u0000\u0000\u00b5\u00b7\u0001\u0000\u0000\u0000"+
		"\u00b6\u00b4\u0001\u0000\u0000\u0000\u00b7\u00b8\u0005\u0004\u0000\u0000"+
		"\u00b8\u0019\u0001\u0000\u0000\u0000\u00b9\u00ba\u0005\u0012\u0000\u0000"+
		"\u00ba\u00bb\u0005%\u0000\u0000\u00bb\u00bc\u0005\u0005\u0000\u0000\u00bc"+
		"\u00bd\u0003\u001e\u000f\u0000\u00bd\u001b\u0001\u0000\u0000\u0000\u00be"+
		"\u00c3\u0005\u0012\u0000\u0000\u00bf\u00c0\u0005\u0012\u0000\u0000\u00c0"+
		"\u00c1\u0005\u0017\u0000\u0000\u00c1\u00c3\u0005\u0012\u0000\u0000\u00c2"+
		"\u00be\u0001\u0000\u0000\u0000\u00c2\u00bf\u0001\u0000\u0000\u0000\u00c3"+
		"\u001d\u0001\u0000\u0000\u0000\u00c4\u00ca\u0005\u0003\u0000\u0000\u00c5"+
		"\u00c9\u0003 \u0010\u0000\u00c6\u00c9\u0003(\u0014\u0000\u00c7\u00c9\u0003"+
		"V+\u0000\u00c8\u00c5\u0001\u0000\u0000\u0000\u00c8\u00c6\u0001\u0000\u0000"+
		"\u0000\u00c8\u00c7\u0001\u0000\u0000\u0000\u00c9\u00cc\u0001\u0000\u0000"+
		"\u0000\u00ca\u00c8\u0001\u0000\u0000\u0000\u00ca\u00cb\u0001\u0000\u0000"+
		"\u0000\u00cb\u00cd\u0001\u0000\u0000\u0000\u00cc\u00ca\u0001\u0000\u0000"+
		"\u0000\u00cd\u00ce\u0005\u0004\u0000\u0000\u00ce\u001f\u0001\u0000\u0000"+
		"\u0000\u00cf\u00d0\u0005\u0012\u0000\u0000\u00d0\u00d1\u0007\u0006\u0000"+
		"\u0000\u00d1\u00d2\u0003:\u001d\u0000\u00d2!\u0001\u0000\u0000\u0000\u00d3"+
		"\u00d9\u0005\u0012\u0000\u0000\u00d4\u00d6\u0005\u0006\u0000\u0000\u00d5"+
		"\u00d7\u0003$\u0012\u0000\u00d6\u00d5\u0001\u0000\u0000\u0000\u00d6\u00d7"+
		"\u0001\u0000\u0000\u0000\u00d7\u00d8\u0001\u0000\u0000\u0000\u00d8\u00da"+
		"\u0005\u0007\u0000\u0000\u00d9\u00d4\u0001\u0000\u0000\u0000\u00d9\u00da"+
		"\u0001\u0000\u0000\u0000\u00da\u00dd\u0001\u0000\u0000\u0000\u00db\u00dc"+
		"\u0005(\u0000\u0000\u00dc\u00de\u0003P(\u0000\u00dd\u00db\u0001\u0000"+
		"\u0000\u0000\u00dd\u00de\u0001\u0000\u0000\u0000\u00de\u00df\u0001\u0000"+
		"\u0000\u0000\u00df\u00e3\u0005\u0003\u0000\u0000\u00e0\u00e2\u0003(\u0014"+
		"\u0000\u00e1\u00e0\u0001\u0000\u0000\u0000\u00e2\u00e5\u0001\u0000\u0000"+
		"\u0000\u00e3\u00e1\u0001\u0000\u0000\u0000\u00e3\u00e4\u0001\u0000\u0000"+
		"\u0000\u00e4\u00e6\u0001\u0000\u0000\u0000\u00e5\u00e3\u0001\u0000\u0000"+
		"\u0000\u00e6\u00e7\u0005\u0004\u0000\u0000\u00e7#\u0001\u0000\u0000\u0000"+
		"\u00e8\u00ed\u0003&\u0013\u0000\u00e9\u00ea\u0005\b\u0000\u0000\u00ea"+
		"\u00ec\u0003&\u0013\u0000\u00eb\u00e9\u0001\u0000\u0000\u0000\u00ec\u00ef"+
		"\u0001\u0000\u0000\u0000\u00ed\u00eb\u0001\u0000\u0000\u0000\u00ed\u00ee"+
		"\u0001\u0000\u0000\u0000\u00ee%\u0001\u0000\u0000\u0000\u00ef\u00ed\u0001"+
		"\u0000\u0000\u0000\u00f0\u00f1\u0005\u0012\u0000\u0000\u00f1\u00f2\u0005"+
		"\u0005\u0000\u0000\u00f2\u00f3\u0003P(\u0000\u00f3\'\u0001\u0000\u0000"+
		"\u0000\u00f4\u00fc\u0003*\u0015\u0000\u00f5\u00fc\u0003,\u0016\u0000\u00f6"+
		"\u00fc\u0003.\u0017\u0000\u00f7\u00fc\u00030\u0018\u0000\u00f8\u00fc\u0003"+
		"2\u0019\u0000\u00f9\u00fc\u00034\u001a\u0000\u00fa\u00fc\u00036\u001b"+
		"\u0000\u00fb\u00f4\u0001\u0000\u0000\u0000\u00fb\u00f5\u0001\u0000\u0000"+
		"\u0000\u00fb\u00f6\u0001\u0000\u0000\u0000\u00fb\u00f7\u0001\u0000\u0000"+
		"\u0000\u00fb\u00f8\u0001\u0000\u0000\u0000\u00fb\u00f9\u0001\u0000\u0000"+
		"\u0000\u00fb\u00fa\u0001\u0000\u0000\u0000\u00fc)\u0001\u0000\u0000\u0000"+
		"\u00fd\u00fe\u0007\u0007\u0000\u0000\u00fe\u00ff\u0003:\u001d\u0000\u00ff"+
		"+\u0001\u0000\u0000\u0000\u0100\u0101\u0007\b\u0000\u0000\u0101\u0102"+
		"\u0005\u0012\u0000\u0000\u0102\u0103\u0005*\u0000\u0000\u0103\u0104\u0003"+
		":\u001d\u0000\u0104-\u0001\u0000\u0000\u0000\u0105\u0106\u0005+\u0000"+
		"\u0000\u0106\u0107\u0003:\u001d\u0000\u0107\u0108\u0005,\u0000\u0000\u0108"+
		"\u010b\u0003\u001e\u000f\u0000\u0109\u010a\u0005-\u0000\u0000\u010a\u010c"+
		"\u0003\u001e\u000f\u0000\u010b\u0109\u0001\u0000\u0000\u0000\u010b\u010c"+
		"\u0001\u0000\u0000\u0000\u010c/\u0001\u0000\u0000\u0000\u010d\u0111\u0005"+
		".\u0000\u0000\u010e\u010f\u0005/\u0000\u0000\u010f\u0111\u0003:\u001d"+
		"\u0000\u0110\u010d\u0001\u0000\u0000\u0000\u0110\u010e\u0001\u0000\u0000"+
		"\u0000\u0111\u0112\u0001\u0000\u0000\u0000\u0112\u0113\u0003\u001e\u000f"+
		"\u0000\u01131\u0001\u0000\u0000\u0000\u0114\u0115\u00050\u0000\u0000\u0115"+
		"\u0116\u0005\u0012\u0000\u0000\u0116\u0117\u0005*\u0000\u0000\u0117\u011a"+
		"\u0005\u0012\u0000\u0000\u0118\u0119\u0005\u0017\u0000\u0000\u0119\u011b"+
		"\u00038\u001c\u0000\u011a\u0118\u0001\u0000\u0000\u0000\u011a\u011b\u0001"+
		"\u0000\u0000\u0000\u011b3\u0001\u0000\u0000\u0000\u011c\u011d\u00051\u0000"+
		"\u0000\u011d\u011f\u0003:\u001d\u0000\u011e\u0120\u0007\t\u0000\u0000"+
		"\u011f\u011e\u0001\u0000\u0000\u0000\u011f\u0120\u0001\u0000\u0000\u0000"+
		"\u01205\u0001\u0000\u0000\u0000\u0121\u0122\u00052\u0000\u0000\u0122\u0123"+
		"\u0003:\u001d\u0000\u01237\u0001\u0000\u0000\u0000\u0124\u0129\u0003:"+
		"\u001d\u0000\u0125\u0126\u0005\b\u0000\u0000\u0126\u0128\u0003:\u001d"+
		"\u0000\u0127\u0125\u0001\u0000\u0000\u0000\u0128\u012b\u0001\u0000\u0000"+
		"\u0000\u0129\u0127\u0001\u0000\u0000\u0000\u0129\u012a\u0001\u0000\u0000"+
		"\u0000\u012a9\u0001\u0000\u0000\u0000\u012b\u0129\u0001\u0000\u0000\u0000"+
		"\u012c\u012d\u0003<\u001e\u0000\u012d;\u0001\u0000\u0000\u0000\u012e\u0133"+
		"\u0003>\u001f\u0000\u012f\u0130\u00053\u0000\u0000\u0130\u0132\u0003>"+
		"\u001f\u0000\u0131\u012f\u0001\u0000\u0000\u0000\u0132\u0135\u0001\u0000"+
		"\u0000\u0000\u0133\u0131\u0001\u0000\u0000\u0000\u0133\u0134\u0001\u0000"+
		"\u0000\u0000\u0134=\u0001\u0000\u0000\u0000\u0135\u0133\u0001\u0000\u0000"+
		"\u0000\u0136\u013b\u0003@ \u0000\u0137\u0138\u00054\u0000\u0000\u0138"+
		"\u013a\u0003@ \u0000\u0139\u0137\u0001\u0000\u0000\u0000\u013a\u013d\u0001"+
		"\u0000\u0000\u0000\u013b\u0139\u0001\u0000\u0000\u0000\u013b\u013c\u0001"+
		"\u0000\u0000\u0000\u013c?\u0001\u0000\u0000\u0000\u013d\u013b\u0001\u0000"+
		"\u0000\u0000\u013e\u0142\u0003B!\u0000\u013f\u0140\u0003X,\u0000\u0140"+
		"\u0141\u0003B!\u0000\u0141\u0143\u0001\u0000\u0000\u0000\u0142\u013f\u0001"+
		"\u0000\u0000\u0000\u0142\u0143\u0001\u0000\u0000\u0000\u0143A\u0001\u0000"+
		"\u0000\u0000\u0144\u014a\u0003D\"\u0000\u0145\u0146\u0003Z-\u0000\u0146"+
		"\u0147\u0003D\"\u0000\u0147\u0149\u0001\u0000\u0000\u0000\u0148\u0145"+
		"\u0001\u0000\u0000\u0000\u0149\u014c\u0001\u0000\u0000\u0000\u014a\u0148"+
		"\u0001\u0000\u0000\u0000\u014a\u014b\u0001\u0000\u0000\u0000\u014bC\u0001"+
		"\u0000\u0000\u0000\u014c\u014a\u0001\u0000\u0000\u0000\u014d\u0153\u0003"+
		"F#\u0000\u014e\u014f\u0003\\.\u0000\u014f\u0150\u0003F#\u0000\u0150\u0152"+
		"\u0001\u0000\u0000\u0000\u0151\u014e\u0001\u0000\u0000\u0000\u0152\u0155"+
		"\u0001\u0000\u0000\u0000\u0153\u0151\u0001\u0000\u0000\u0000\u0153\u0154"+
		"\u0001\u0000\u0000\u0000\u0154E\u0001\u0000\u0000\u0000\u0155\u0153\u0001"+
		"\u0000\u0000\u0000\u0156\u0159\u0003H$\u0000\u0157\u0158\u0005\t\u0000"+
		"\u0000\u0158\u015a\u0003F#\u0000\u0159\u0157\u0001\u0000\u0000\u0000\u0159"+
		"\u015a\u0001\u0000\u0000\u0000\u015aG\u0001\u0000\u0000\u0000\u015b\u015c"+
		"\u0007\n\u0000\u0000\u015c\u015f\u0003J%\u0000\u015d\u015f\u0003J%\u0000"+
		"\u015e\u015b\u0001\u0000\u0000\u0000\u015e\u015d\u0001\u0000\u0000\u0000"+
		"\u015fI\u0001\u0000\u0000\u0000\u0160\u016b\u0005\u0013\u0000\u0000\u0161"+
		"\u016b\u0005\u0014\u0000\u0000\u0162\u016b\u0005\u0015\u0000\u0000\u0163"+
		"\u016b\u0005\u0012\u0000\u0000\u0164\u016b\u0003L&\u0000\u0165\u0166\u0005"+
		"\u0006\u0000\u0000\u0166\u0167\u0003:\u001d\u0000\u0167\u0168\u0005\u0007"+
		"\u0000\u0000\u0168\u016b\u0001\u0000\u0000\u0000\u0169\u016b\u0003N\'"+
		"\u0000\u016a\u0160\u0001\u0000\u0000\u0000\u016a\u0161\u0001\u0000\u0000"+
		"\u0000\u016a\u0162\u0001\u0000\u0000\u0000\u016a\u0163\u0001\u0000\u0000"+
		"\u0000\u016a\u0164\u0001\u0000\u0000\u0000\u016a\u0165\u0001\u0000\u0000"+
		"\u0000\u016a\u0169\u0001\u0000\u0000\u0000\u016bK\u0001\u0000\u0000\u0000"+
		"\u016c\u016d\u0005\u0012\u0000\u0000\u016d\u016f\u0005\u0006\u0000\u0000"+
		"\u016e\u0170\u00038\u001c\u0000\u016f\u016e\u0001\u0000\u0000\u0000\u016f"+
		"\u0170\u0001\u0000\u0000\u0000\u0170\u0171\u0001\u0000\u0000\u0000\u0171"+
		"\u0172\u0005\u0007\u0000\u0000\u0172M\u0001\u0000\u0000\u0000\u0173\u017c"+
		"\u0005\u0003\u0000\u0000\u0174\u0179\u0003:\u001d\u0000\u0175\u0176\u0005"+
		"\b\u0000\u0000\u0176\u0178\u0003:\u001d\u0000\u0177\u0175\u0001\u0000"+
		"\u0000\u0000\u0178\u017b\u0001\u0000\u0000\u0000\u0179\u0177\u0001\u0000"+
		"\u0000\u0000\u0179\u017a\u0001\u0000\u0000\u0000\u017a\u017d\u0001\u0000"+
		"\u0000\u0000\u017b\u0179\u0001\u0000\u0000\u0000\u017c\u0174\u0001\u0000"+
		"\u0000\u0000\u017c\u017d\u0001\u0000\u0000\u0000\u017d\u017e\u0001\u0000"+
		"\u0000\u0000\u017e\u017f\u0005\u0004\u0000\u0000\u017fO\u0001\u0000\u0000"+
		"\u0000\u0180\u0182\u0003R)\u0000\u0181\u0183\u0005\n\u0000\u0000\u0182"+
		"\u0181\u0001\u0000\u0000\u0000\u0182\u0183\u0001\u0000\u0000\u0000\u0183"+
		"Q\u0001\u0000\u0000\u0000\u0184\u0185\u0007\u000b\u0000\u0000\u0185S\u0001"+
		"\u0000\u0000\u0000\u0186\u018b\u0003*\u0015\u0000\u0187\u018b\u0003,\u0016"+
		"\u0000\u0188\u018b\u00034\u001a\u0000\u0189\u018b\u00032\u0019\u0000\u018a"+
		"\u0186\u0001\u0000\u0000\u0000\u018a\u0187\u0001\u0000\u0000\u0000\u018a"+
		"\u0188\u0001\u0000\u0000\u0000\u018a\u0189\u0001\u0000\u0000\u0000\u018b"+
		"U\u0001\u0000\u0000\u0000\u018c\u0190\u0005\u000b\u0000\u0000\u018d\u018f"+
		"\t\u0000\u0000\u0000\u018e\u018d\u0001\u0000\u0000\u0000\u018f\u0192\u0001"+
		"\u0000\u0000\u0000\u0190\u0191\u0001\u0000\u0000\u0000\u0190\u018e\u0001"+
		"\u0000\u0000\u0000\u0191\u0193\u0001\u0000\u0000\u0000\u0192\u0190\u0001"+
		"\u0000\u0000\u0000\u0193\u0194\u0005\u000b\u0000\u0000\u0194W\u0001\u0000"+
		"\u0000\u0000\u0195\u0196\u0007\f\u0000\u0000\u0196Y\u0001\u0000\u0000"+
		"\u0000\u0197\u0198\u0007\r\u0000\u0000\u0198[\u0001\u0000\u0000\u0000"+
		"\u0199\u019a\u0007\u000e\u0000\u0000\u019a]\u0001\u0000\u0000\u0000&a"+
		"l{\u0082\u008a\u0094\u0096\u009d\u00a9\u00b4\u00c2\u00c8\u00ca\u00d6\u00d9"+
		"\u00dd\u00e3\u00ed\u00fb\u010b\u0110\u011a\u011f\u0129\u0133\u013b\u0142"+
		"\u014a\u0153\u0159\u015e\u016a\u016f\u0179\u017c\u0182\u018a\u0190";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}