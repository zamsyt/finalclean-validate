# TFC-validate (2023)

## Instructions

1. Save the Drawpile as a .ora file
2. With Go installed, in the repo root, run `go run . <command> <path-to-drawpile.ora>`

## Commands

### `checkpalette <my-drawpile.ora>`

Count the pixels on each layer that don't match the palette (defined by r-slash-place-2023.gpl). Create an image for each layer representing the deviation of each pixel from the palette. (saved under a "diff" subdirectory)

### `fixpalette <my-drawpile.ora>`

For each layer, create a version with colors converted to match the palette (using an euclidean RGB distance). Already correct images are skipped. (images are saved under a "fix" subdirectory)

#### Flags

- `--alpha-threshold X` number between 0-255. Alpha values lower than this become transparent. Ones that are higher are converted to the closest color from the palette.
- `--diff` only draw the changed pixels
- `--crop` crop transparent edges off the image
- `--fullsize` save each layer at the full canvas size (3000x2000)
