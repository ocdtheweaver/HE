grammar HE;

program: line* EOF;

line: summon
    | withAssets
    | object
    | globalStatement
    | comment
    | NEWLINE
    ;

summon: SUMMON STRING (NAMED | AS) ID;

withAssets: WITH assetList;
assetList: asset (AND asset)*;
asset: assetType STRING (NAMED ID)?;
assetType: IMAGE | SOUND | MUSIC | VIDEO | FONT | SHADER;

object: (CREATE | MAKE) ID (LIKE ID)? '{' objectBody '}';
objectBody: (section | object | comment | withAssets)*;

section: properties | abilities | reactions | memories;

properties: ID (HAS | OWNS | CARRIES) block;
abilities: ID (CAN | KNOWS_HOW_TO) '[' action* ']';
reactions: (WHEN | WHENEVER | ON) trigger '[' statement* ']';
memories: ID REMEMBERS ':' block;

trigger: ID | ID WITH ID;

block: '[' (property | statement | comment)* ']';
property: ID (IS | STARTS_AS) expression;

action: ID ('(' paramList? ')')? (RETURNS type)? '[' statement* ']';
paramList: parameter (',' parameter)*;
parameter: ID ':' type;

statement: say | change | decide | repeat | call | waitStmt | returnStmt;  // CHANGED: waitStmt

say: (SAY | PRINT) expression;
change: (SET | MAKE) ID TO expression;
decide: IF expression THEN block (ELSE block)?;
repeat: (REPEAT | WHILE expression) block;
call: TELL ID TO ID (WITH argList)?;
waitStmt: WAIT expression (SECONDS | FRAMES)?;  // CHANGED: waitStmt
returnStmt: RETURN expression;

argList: expression (',' expression)*;

expression: logicOr;
logicOr: logicAnd (OR logicAnd)*;
logicAnd: comparison (AND comparison)*;
comparison: arithmetic (compOp arithmetic)?;
arithmetic: term (addOp term)*;
term: factor (multOp factor)*;
factor: unary ('^' factor)?;
unary: (MINUS | BANG) primary | primary;
primary: NUMBER 
       | STRING 
       | BOOLEAN 
       | ID 
       | callExpression
       | '(' expression ')' 
       | arrayLiteral
       ;

callExpression: ID '(' argList? ')';
arrayLiteral: '[' (expression (',' expression)*)? ']';

type: baseType ('[]')?;
baseType: NUMBER_TYPE | STRING_TYPE | BOOLEAN_TYPE;

globalStatement: say | change | waitStmt | call;  // CHANGED: waitStmt

comment: '~' .*? '~';

// Lexer rules
ID: [a-zA-Z_][a-zA-Z0-9_]*;
NUMBER: [0-9]+ ('.' [0-9]+)?;
STRING: '"' (~["\\] | '\\' .)* '"';
BOOLEAN: 'true' | 'false' | 'yes' | 'no';

// Operators
compOp: '==' | '!=' | '>' | '<' | '>=' | '<=';
addOp: PLUS | MINUS;
multOp: STAR | SLASH;

// Keywords
SUMMON: 'summon';
WITH: 'with';
NAMED: 'named';
AS: 'as';
CREATE: 'create';
MAKE: 'make';
LIKE: 'like';
HAS: 'has';
OWNS: 'owns';
CARRIES: 'carries';
CAN: 'can';
KNOWS_HOW_TO: 'knows' 'how' 'to';
WHEN: 'when';
WHENEVER: 'whenever';
ON: 'on';
REMEMBERS: 'remembers';
IS: 'is';
STARTS_AS: 'starts as';
RETURNS: 'returns';
SET: 'set';
TO: 'to';
IF: 'if';
THEN: 'then';
ELSE: 'else';
REPEAT: 'repeat';
WHILE: 'while';
TELL: 'tell';
WAIT: 'wait';
RETURN: 'return';
OR: 'or';
AND: 'and';
NOT: 'not';

// Asset types
IMAGE: 'image';
SOUND: 'sound';
MUSIC: 'music';
VIDEO: 'video';
FONT: 'font';
SHADER: 'shader';

// Statement keywords
SAY: 'say';
PRINT: 'print';
SECONDS: 'seconds';
FRAMES: 'frames';

// Types
NUMBER_TYPE: 'number';
STRING_TYPE: 'string';
BOOLEAN_TYPE: 'boolean';

// Symbols
PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
BANG: '!';

// Whitespace
NEWLINE: '\r'? '\n';
WS: [ \t]+ -> skip;