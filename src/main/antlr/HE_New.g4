grammar HE_New;

@header {
package he_new.parser;
}

program: line* EOF;

line: summon
    | withAssets
    | object
    | NEWLINE
    ;

summon: SUMMON STRING (AS | NAMED) ID;

withAssets: WITH assetList;
assetList: asset (AND asset)*;            // AND is defined below
asset: assetType STRING (NAMED ID)?;
assetType: IMAGE | SOUND | MUSIC | VIDEO | FONT | SHADER;

object: MAKE ID (LIKE ID)? LBRACE objectBody RBRACE;
objectBody: objectSection*;

objectSection: propertiesSection | abilitiesSection | reactionSection;

propertiesSection: ID HAS COLON propertyBlock+;
propertyBlock: LBRACK property RBRACK;
property: ID IS expression;

abilitiesSection: ID CAN COLON abilityBlock+;
abilityBlock: method LBRACK statement* RBRACK;
method: ID LPAREN paramList? RPAREN;
paramList: parameter (COMMA parameter)*;
parameter: ID COLON type;

reactionSection: ON ID LBRACK statement* RBRACK;

type: STRING_TYPE | NUMBER_TYPE | BOOLEAN_TYPE;

statement: 
    printStmt
    | callStmt
    | showStmt
    | waitStmt
    | assignmentStmt
    ;

printStmt: PRINT expression;
callStmt: TELL ID TO ID (WITH argList)?;
showStmt: SHOW ID CENTERED?;
waitStmt: WAIT expression SECONDS;
assignmentStmt: ID IS expression;

argList: expression (COMMA expression)*;

// FIXED: Non-left-recursive expression grammar
expression: additiveExpression;

additiveExpression: 
    multiplicativeExpression ((PLUS | MINUS) multiplicativeExpression)*
    ;

multiplicativeExpression: 
    primaryExpression ((MULT | DIV) primaryExpression)*
    ;

primaryExpression: 
    STRING
    | ID
    | NUMBER
    | LPAREN expression RPAREN
    ;

// ========== LEXER RULES ==========
// Keywords
SUMMON: 'summon';
WITH: 'with';
NAMED: 'named';
AS: 'as';
MAKE: 'make';
LIKE: 'like';
HAS: 'has';
CAN: 'can';
ON: 'on';
IS: 'is';
PRINT: 'print';
TELL: 'tell';
TO: 'to';
SHOW: 'show';
CENTERED: 'centered';
WAIT: 'wait';
SECONDS: 'seconds';

// FIXED: Added missing AND token
AND: 'and';

// Asset types
IMAGE: 'image';
SOUND: 'sound';
MUSIC: 'music';
VIDEO: 'video';
FONT: 'font';
SHADER: 'shader';

// Types
STRING_TYPE: 'string';
NUMBER_TYPE: 'number';
BOOLEAN_TYPE: 'boolean';

// Symbols
COLON: ':';
LBRACK: '[';
RBRACK: ']';
LPAREN: '(';
RPAREN: ')';
COMMA: ',';
TILDA: '~';
LBRACE: '{';
RBRACE: '}';
DOT: '.';
PLUS: '+';
MINUS: '-';
MULT: '*';
DIV: '/';

// Comments go to HIDDEN channel
COMMENT: TILDA .*? TILDA -> channel(HIDDEN);

// Basic tokens
ID: [a-zA-Z_][a-zA-Z0-9_]* (DOT [a-zA-Z_][a-zA-Z0-9_]*)*;
NUMBER: [0-9]+ ('.' [0-9]+)?;
STRING: '"' (~["\\] | '\\' .)* '"';

// Whitespace
NEWLINE: '\r'? '\n';
WS: [ \t]+ -> skip;