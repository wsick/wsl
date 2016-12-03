# WSL 

(Wicked Sick Language pronounced 'wassail')

Domain-specific language intended to represent user interfaces.

## Objectives

The purpose of this language is to take the great parts of XML, JSON, and HCL with user interfaces in mind.  
The language should be dense, terse, and easy to parse.  
This will allow for faster transmission, easier reading on the screen, and improved parse performance.

Notably, the following considerations are made:
- 1. Within declaration context `(...)`, whitespace after a value indicates new declaration
- 2. Within property context `[...]`, whitespace after a value indicates new property
- 3. Within body context `{...}`, whitespace after a full element indicates new element
- 4. Element and Property names are lowercased. (`-` is used to break words apart)

Considerations 1-3 enable the ability to create streaming ASTs with no backtracking.
Consideration 4 enables simplistic conversion between XAML and WSL. 

## Examples

See example transpilations to wsl:
  - Example 1:
    - Input: [inputs/toy1.xaml](inputs/toy1.xaml)
    - Output: [outputs/toy1.wsl](outputs/toy1.wsl)
 
