import { expect, test } from "vitest"
import {
    fs,
} from "./challenge.js"

test("test function-state", () => {
    expect(fs()).toEqual(42);
    expect(fs()).toEqual(420);
    expect(fs()).toEqual(69);
    expect(fs()).toEqual(69);
    expect(fs()).toEqual(420);
    expect(fs()).toEqual(42);
});



