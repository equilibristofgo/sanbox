#include "libmylib.h"

int main(void) {
    say(const_cast<char*>("hello world"));

    return 0;
}