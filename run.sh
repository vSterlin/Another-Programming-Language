go run .;
llvm-gcc ./build/out.ll -o ./build/app;
./build/app; 
echo $?;