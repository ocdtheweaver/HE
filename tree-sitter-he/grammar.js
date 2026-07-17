/**
 * tree-sitter-he — Tree-sitter grammar for the HE programming language.
 *
 * HE is a near-English programming language. Its key syntactic features:
 *   - Tilde comments:  ~ this is a comment ~
 *   - Blocks use [ and ] instead of { }
 *   - String interpolation with { } inside strings: "Hello, {name}!"
 *   - #protected[N] protection tags on objects/abilities
 *   - Keywords read like English: say, set, tell, create, make, summon, etc.
 *
 * Grammar version: HE v5.0.0 (Pass 14)
 */

module.exports = grammar({
  name: 'he',

  // Characters that can appear in the middle of tokens
  extras: $ => [
    /\s/,
    $.comment,
  ],

  // When the parser sees ambiguity, prefer these rules
  conflicts: $ => [
    [$.call_stmt, $.expr_stmt],
    [$.method_call, $.field_access],
    [$.set_stmt, $.multi_assign_stmt],
  ],

  // Words that can't be used as identifiers
  word: $ => $.identifier,

  rules: {

    // ── Top-level ───────────────────────────────────────────────────────────
    source_file: $ => repeat($._line),

    _line: $ => choice(
      $.summon_stmt,
      $.object_def,
      $.statement,
    ),

    // ── Comments ────────────────────────────────────────────────────────────
    // ~ anything here, can span multiple lines ~
    comment: $ => token(seq('~', /[^~]*/, '~')),

    // ── Identifiers and literals ─────────────────────────────────────────────
    identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,

    number: $ => /[0-9]+(\.[0-9]+)?/,

    boolean: $ => choice('true', 'yes', 'false', 'no'),

    nothing: $ => 'nothing',

    // Plain string (no interpolation)
    string: $ => seq(
      '"',
      repeat(choice(
        /[^"{}\\]+/,
        /\\./,             // escape sequences
      )),
      '"',
    ),

    // Interpolated string: "Hello, {name}! You have {count} items."
    // The {…} regions contain full expressions.
    interp_string: $ => seq(
      '"',
      repeat(choice(
        /[^"{}\\]+/,       // plain text
        /\\./,             // escape sequences
        seq(
          '{',
          $._expression,   // full expression inside { }
          '}',
        ),
      )),
      '"',
    ),

    // Protection tag: #protected, #protected1, #protected2, ...
    // Lexed as '#' immediately followed by an identifier (no space).
    protect_tag: $ => token(seq('#', /[a-zA-Z][a-zA-Z0-9]*/)),

    // ── Summon (import) ──────────────────────────────────────────────────────
    summon_stmt: $ => seq(
      'summon',
      field('module', choice($.string, $.interp_string)),
      optional(seq(
        choice('as', 'named'),
        field('alias', $.identifier),
      )),
    ),

    // ── Object definitions ───────────────────────────────────────────────────
    object_def: $ => seq(
      field('kind', choice('create', 'make')),
      field('name', $.identifier),
      optional(field('protect_tag', $.protect_tag)),
      optional(seq('like', field('parent', $.identifier))),
      '[',
      repeat($._object_section),
      ']',
    ),

    _object_section: $ => choice(
      $.has_section,
      $.can_section,
      $.on_section,
    ),

    // has: [ name is "value", health is 100 ]
    has_section: $ => seq(
      choice('has', 'owns', 'carries', 'remembers'),
      ':',
      '[',
      repeat($.property),
      ']',
    ),

    property: $ => seq(
      field('name', $.identifier),
      optional(seq(':', field('type_hint', $.identifier))),  // optional type annotation
      choice('is', 'becomes', seq('starts', choice('as', 'be'))),
      field('value', $._expression),
    ),

    // can: [ methodName(params) #tag [ body ] ]
    can_section: $ => seq(
      choice('can', seq('know', 'how', 'to')),
      ':',
      '[',
      repeat($.ability_def),
      ']',
    ),

    ability_def: $ => seq(
      field('name', $.identifier),
      optional(field('protect_tag', $.protect_tag)),   // before params
      optional($.param_list),
      optional(field('protect_tag_after', $.protect_tag)),  // or after params
      optional(seq(
        choice('returns', 'return'),
        field('return_type', $.identifier),
      )),
      '[',
      repeat($.statement),
      ']',
    ),

    param_list: $ => seq(
      '(',
      optional(seq(
        $.param,
        repeat(seq(',', $.param)),
      )),
      ')',
    ),

    param: $ => seq(
      field('name', $.identifier),
      optional(seq(':', field('type', $.identifier))),
    ),

    // on click [ body ] / on collision [ body ]
    on_section: $ => seq(
      choice('on', 'when', 'whenever'),
      field('event', $.identifier),
      optional(seq('with', field('value', $._expression))),
      '[',
      repeat($.statement),
      ']',
    ),

    // ── Statements ────────────────────────────────────────────────────────────
    statement: $ => choice(
      $.say_stmt,
      $.set_stmt,
      $.multi_assign_stmt,
      $.becomes_stmt,
      $.change_stmt,
      $.grow_stmt,
      $.shrink_stmt,
      $.dot_assign_stmt,
      $.give_stmt,
      $.tell_stmt,
      $.if_stmt,
      $.check_if_stmt,
      $.repeat_while_stmt,
      $.repeat_times_stmt,
      $.counted_repeat_stmt,
      $.repeat_until_stmt,
      $.for_each_stmt,
      $.for_each_range_stmt,
      $.each_range_stmt,
      $.for_each_field_stmt,
      $.try_stmt,
      $.try_with_var_stmt,
      $.return_stmt,
      $.multi_return_stmt,
      $.remember_stmt,
      $.recall_stmt,
      $.forget_stmt,
      $.wait_stmt,
      $.ask_stmt,
      $.with_scope_stmt,
      $.summon_stmt,     // summon inside a method body
      $.expr_stmt,
    ),

    // say / print / show "…"
    say_stmt: $ => seq(
      choice('say', 'print', 'show'),
      $._expression,
    ),

    // set x to expr
    set_stmt: $ => seq(
      choice('set', 'let'),
      field('name', $.identifier),
      optional(seq(
        '.',
        field('field', $.identifier),
      )),
      choice('to', 'be', 'is', 'becomes'),
      field('value', $._expression),
    ),

    // set x, y to a, b
    multi_assign_stmt: $ => seq(
      choice('set', 'let'),
      field('names', seq(
        $.identifier,
        repeat1(seq(',', $.identifier)),
      )),
      choice('to', 'be'),
      field('values', seq(
        $._expression,
        repeat(seq(',', $._expression)),
      )),
    ),

    // name becomes value  (standalone mutation)
    becomes_stmt: $ => seq(
      field('name', $.identifier),
      'becomes',
      field('value', $._expression),
    ),

    // change x to expr
    change_stmt: $ => seq(
      'change',
      field('name', $.identifier),
      choice('to', 'be'),
      field('value', $._expression),
    ),

    // grow x by n
    grow_stmt: $ => seq(
      'grow',
      field('name', $.identifier),
      'by',
      field('amount', $._expression),
    ),

    // shrink x by n
    shrink_stmt: $ => seq(
      'shrink',
      field('name', $.identifier),
      'by',
      field('amount', $._expression),
    ),

    // set obj.field to expr
    dot_assign_stmt: $ => seq(
      choice('set', 'change'),
      field('object', $.identifier),
      '.',
      field('field', $.identifier),
      choice('to', 'be', 'is'),
      field('value', $._expression),
    ),

    // give expr to obj.field
    give_stmt: $ => seq(
      'give',
      field('value', $._expression),
      'to',
      field('object', $.identifier),
      '.',
      field('field', $.identifier),
    ),

    // tell X to Y [with args]
    tell_stmt: $ => seq(
      'tell',
      field('object', $.identifier),
      'to',
      field('action', $.identifier),
      optional(seq(
        'with',
        field('args', $._expression),
        repeat(seq(',', $._expression)),
      )),
    ),

    // if cond then [ body ] [else [ body ]]
    if_stmt: $ => seq(
      'if',
      field('condition', $._expression),
      'then',
      '[',
      repeat($.statement),
      ']',
      optional(seq(
        'else',
        '[',
        repeat($.statement),
        ']',
      )),
    ),

    // check if cond then [...]  (alias for if)
    check_if_stmt: $ => seq(
      'check',
      'if',
      field('condition', $._expression),
      'then',
      '[',
      repeat($.statement),
      ']',
      optional(seq(
        'else',
        '[',
        repeat($.statement),
        ']',
      )),
    ),

    // repeat while cond [...]
    repeat_while_stmt: $ => seq(
      'repeat',
      'while',
      field('condition', $._expression),
      '[',
      repeat($.statement),
      ']',
    ),

    // repeat N times [...]
    repeat_times_stmt: $ => seq(
      'repeat',
      field('count', $._expression),
      'times',
      '[',
      repeat($.statement),
      ']',
    ),

    // repeat N times as i [...]  (counter exposed)
    counted_repeat_stmt: $ => seq(
      'repeat',
      field('count', $._expression),
      'times',
      'as',
      field('counter', $.identifier),
      '[',
      repeat($.statement),
      ']',
    ),

    // repeat until cond [...]
    repeat_until_stmt: $ => seq(
      'repeat',
      'until',
      field('condition', $._expression),
      '[',
      repeat($.statement),
      ']',
    ),

    // for each item in list [...]
    for_each_stmt: $ => seq(
      'for',
      'each',
      field('var', $.identifier),
      'in',
      field('list', $._expression),
      '[',
      repeat($.statement),
      ']',
    ),

    // for each i from A to B [step S] [...]
    for_each_range_stmt: $ => seq(
      'for',
      'each',
      field('var', $.identifier),
      'from',
      field('from', $._expression),
      choice('to', 'upto'),
      field('to', $._expression),
      optional(seq('step', field('step', $._expression))),
      '[',
      repeat($.statement),
      ']',
    ),

    // each i from A to B [step S] [...]  (short form, no 'for')
    each_range_stmt: $ => seq(
      'each',
      field('var', $.identifier),
      'from',
      field('from', $._expression),
      choice('to', 'upto'),
      field('to', $._expression),
      optional(seq('step', field('step', $._expression))),
      '[',
      repeat($.statement),
      ']',
    ),

    // for each key, val in obj [...]
    for_each_field_stmt: $ => seq(
      'for',
      'each',
      field('key', $.identifier),
      ',',
      field('val', $.identifier),
      'in',
      field('object', $._expression),
      '[',
      repeat($.statement),
      ']',
    ),

    // try [ body ] or [ handler ]
    try_stmt: $ => seq(
      'try',
      '[',
      repeat($.statement),
      ']',
      'or',
      optional(seq('if', 'it', 'fails')),   // "or if it fails [...]"
      '[',
      repeat($.statement),
      ']',
    ),

    // try [ body ] or (errVar) [ handler ]
    try_with_var_stmt: $ => seq(
      'try',
      '[',
      repeat($.statement),
      ']',
      'or',
      '(',
      field('error_var', $.identifier),
      ')',
      '[',
      repeat($.statement),
      ']',
    ),

    // return expr
    return_stmt: $ => seq(
      'return',
      $._expression,
    ),

    // return a, b, c
    multi_return_stmt: $ => seq(
      'return',
      $._expression,
      repeat1(seq(',', $._expression)),
    ),

    // remember x [as "key"]
    remember_stmt: $ => seq(
      'remember',
      field('name', $.identifier),
      optional(seq('as', field('key', $.string))),
    ),

    // recall x [as alias]
    recall_stmt: $ => seq(
      'recall',
      field('name', $.identifier),
      optional(seq('as', field('alias', $.identifier))),
    ),

    // forget x
    forget_stmt: $ => seq(
      'forget',
      field('name', $.identifier),
    ),

    // wait N seconds / wait N frames
    wait_stmt: $ => seq(
      'wait',
      field('amount', $._expression),
      choice('seconds', 'frames', 'second', 'frame'),
    ),

    // ask "prompt?" as varName
    ask_stmt: $ => seq(
      'ask',
      field('prompt', $._expression),
      'as',
      field('var', $.identifier),
    ),

    // with expr as name [ body ]
    with_scope_stmt: $ => seq(
      'with',
      field('expr', $._expression),
      'as',
      field('alias', $.identifier),
      '[',
      repeat($.statement),
      ']',
    ),

    // Expression used as a statement (method calls, etc.)
    expr_stmt: $ => $._expression,

    // ── Expressions ───────────────────────────────────────────────────────────
    _expression: $ => choice(
      $.logic_or,
    ),

    logic_or: $ => choice(
      $.logic_and,
      seq($.logic_or, 'or', $.logic_and),
    ),

    logic_and: $ => choice(
      $.comparison,
      seq($.logic_and, 'and', $.comparison),
    ),

    comparison: $ => choice(
      $.arithmetic,
      seq($.arithmetic, '==', $.arithmetic),
      seq($.arithmetic, '!=', $.arithmetic),
      seq($.arithmetic, '>', $.arithmetic),
      seq($.arithmetic, '<', $.arithmetic),
      seq($.arithmetic, '>=', $.arithmetic),
      seq($.arithmetic, '<=', $.arithmetic),
      seq($.arithmetic, 'is', 'not', $.arithmetic),
      seq($.arithmetic, 'is', $.arithmetic),
      // X is between A and B
      seq($.arithmetic, 'is', 'between', $.arithmetic, 'and', $.arithmetic),
      // X is one of [list]
      seq($.arithmetic, 'is', 'one', 'of', $._expression),
      // not expr
      seq('not', $._expression),
    ),

    arithmetic: $ => choice(
      $.term,
      seq($.arithmetic, '+', $.term),
      seq($.arithmetic, '-', $.term),
    ),

    term: $ => choice(
      $.factor,
      seq($.term, '*', $.factor),
      seq($.term, '/', $.factor),
    ),

    factor: $ => choice(
      $.unary,
      seq($.unary, '**', $.factor),  // right-associative power
    ),

    unary: $ => choice(
      $.postfix,
      seq('-', $.postfix),
      seq('not', $.postfix),
      seq('!', $.postfix),
    ),

    postfix: $ => choice(
      $.primary,
      $.method_call,
      $.field_access,
      $.call_expr,
    ),

    // obj.method(args)
    method_call: $ => seq(
      field('receiver', $.postfix),
      '.',
      field('method', $.identifier),
      '(',
      optional($.arg_list),
      ')',
    ),

    // obj.field
    field_access: $ => seq(
      field('receiver', $.postfix),
      '.',
      field('field', $.identifier),
    ),

    // func(args)
    call_expr: $ => seq(
      field('callee', $.identifier),
      '(',
      optional($.arg_list),
      ')',
    ),

    arg_list: $ => seq(
      $._expression,
      repeat(seq(',', $._expression)),
    ),

    // ── Primary expressions ──────────────────────────────────────────────────
    primary: $ => choice(
      $.number,
      $.boolean,
      $.nothing,
      $.interp_string,
      $.string,
      $.array,
      $.named_arg_block,
      $.ability_lit,
      seq('(', $._expression, ')'),
      $.identifier,
    ),

    // [a, b, c]
    array: $ => seq(
      '[',
      optional(seq(
        $._expression,
        repeat(seq(',', $._expression)),
      )),
      ']',
    ),

    // [key is value, key2 is value2]  (named arg block for UI/wolfhead)
    named_arg_block: $ => seq(
      '[',
      seq(
        $.named_pair,
        repeat(seq(',', $.named_pair)),
      ),
      ']',
    ),

    named_pair: $ => seq(
      field('key', $.identifier),
      choice('is', ':', 'becomes'),
      field('value', $._expression),
    ),

    // ability(params) [ body ]
    ability_lit: $ => seq(
      'ability',
      optional(seq(
        '(',
        optional(seq(
          $.identifier,
          repeat(seq(',', $.identifier)),
        )),
        ')',
      )),
      '[',
      repeat($.statement),
      ']',
    ),
  },
});
