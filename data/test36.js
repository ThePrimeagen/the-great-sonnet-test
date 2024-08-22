import { expect, test } from "vitest"
import {
    number2,
} from "./challenge.js"

test("test number2", () => {
    expect(number2.foo).toEqual("7")
    expect(number2().foo + +number2.foo).toEqual(69)
    expect(number2() + number2()).toEqual("nini")
});


