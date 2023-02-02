package main

// TODO:
// * reduce the number of encoding collision keywords
//
// Build 2.
// -- Added the encoding collision sfx
// -- Added the encoding collision manual entry
// -- Renamed "Access Keywords" to "Encoded Keywords"
// -- Fixed ctrl modifier in wasm builds
//
// Build 3.
// > Bug fixes:
// -- Fixed io logs text alignment
// -- Fixed a compound keyword encoded collision bug
// -- Fixed invisible cursor at 10 input letters
// > Content modifications:
// -- Changed the final compound keyword to "cloudburst"
// -- Rebalanced some levels
// -- Fixed a few typos (in various places!)
// -- Replaced some slang to be more player-friendly
// -- Changed the decipher screen music
// > New content:
// -- Added a polygraphic ciphers manual page
// -- Added a conditional transformations manual page
// -- Added a couple of new levels
// -- Added menu music
//
// Build 4.
// > Other:
// -- Made sound levels configurable
// > Content modifications:
// -- Rebalanced some levels
// > New content:
// -- Added a conditions vocab manual page
// -- Added a "binary tree" level
// -- Added a "conveyor" level
//
// Build 5.
// > Bug fixes:
// -- Second checkmark near the bonus levels (if they were completed and moved from the story levels)
// > Content modifications:
// -- A "switch" level now modifies words longer than 5 too
// -- Done some more manual pages proof reading
//
// Build 6.
// > Gameplay:
// -- Increased the game speed scale
// -- Encoding collisions results in red output text
// -- Encoding success results in gold output text
// -- "go to manual" action after clearing the level
// -- Added help stickers to some starting levels
// -- New "custom simulations" mode
// > Content modifications:
// -- Some balance changes
//
// Build 7:
// > UX:
// -- Can now type in letters while holding "shift" (also fixes capslock issue)
// -- Use 'keyname' notation instead [keyname] inside results screen (more readable)
// > Other:
// -- Improved custom levels selection screen
const buildVersion = "7"
