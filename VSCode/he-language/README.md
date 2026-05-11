# HE Language Support

This extension adds syntax highlighting and basic language features for the **HE programming language**.

## Features
- Syntax highlighting
- Block and action detection
- Object definitions
- Loop keywords
- Conditional keywords
- Type highlighting
- Auto-closing brackets
- Indentation rules

## File Extension
`.he`

## Example

```he
fetch physics

make Player {
    
    state [
        speed: dec = 4.5
        jumping: int = 0
        end alive
    ]

    movement() [
        if (jumping == 1)
            print "Player is jumping"

        repeat [
            reduce speed by 1
        ]
        until alive
    ]
}
