# rotate_exercise
Coding Exercise in Rotation

# Goal

Produce a cli application which will take 3 arguments - "left/right", input file, output file.

The application must shift all bytes within the input file one to the left/right in accordance with the input argument, wrapping the final shift to the other end of the file.

# Considerations

1. As the size of the input is unknown, we cannot assume we can load the entire file into memory
