package prompt

import "strings"

type PromptParams struct {
    Export string
    Language string
    TestFile string
    TestRunner string
    Error string
    Code string
    Reasoning string
    Motivation string
}

func CreatePromptParamsFromLanguage(lang string) PromptParams {
    switch (lang) {
    case "JavaScript":
        return PromptParams{
            Export: JS_EXPORT,
            Motivation: "",
            Language: lang,
            TestFile: "",
            TestRunner: "vitest",
            Error: "",
            Code: "",
            Reasoning: "",
        }
    case "Golang":
        return PromptParams{
            Motivation: "",
            Export: "",
            Language: lang,
            TestFile: "",
            TestRunner: "go's built in test runner",
            Error: "",
            Code: "",
            Reasoning: "",
        }
    }

    panic("incorrect language")
}

const JS_EXPORT = "All functions used in tests are imported without default imports"

const TWITCH_PROMPT = `You are an expert Software Engineer delivering code of truly perfect quality.
Your goal is to make all of the unit tests pass by providing the __LANGUAGE__ code required to resolve the errors.

You'll be provided with the following input:
- The output from __TEST_RUNNER__ as error messages.

You're output must be the following:
- A separate file serving as a module that can be run in an isolated environment.
- The output must be stripped of all markdown

Take the following into account when generating the code:
- __EXPORT__
- You never generate poetry, only run-able, modularized __LANGUAGE__ code

Format your output as:
- Pure, valid __LANGUAGE__
- No markdown
- No JSON
- No XML

Example input:
Here is the test file's exact contents:
__TEST_FILE__
`

const MY_PROMPT = `You are a Staff Level Software Engineer with an incredible stock package
and your goal is to make the unit tests pass by providing the __LANGUAGE__ code based
on the errors you receive.

The errors provided in the user prompt will be the output from __TEST_RUNNER__

Your output will be used as a separate file, so make sure you create it as a module

Strip all markdown

__EXPORT__

No poetry please

format your output as valid __LANGUAGE__ ONLY.  Do not provide anything
but valid __LANGUAGE__ code, NO MARKDOWN, NO JSON, NO XML, JUST CODE

DO NOT HALLUCINATE

YOU SHALL WRITE ALL OF THE CODE

Here is the test file's exact contents:
__TEST_FILE__
`

const CODE_GEN_PROMPT = `
<MustReadCarefully>
__MOTIVATION__
</MustReadCarefully>

Take your time to generate the code and think step by step in resolving the issues.
<Reasoning For Error>
__REASONING__
</Reasoning For Error>

<Error>
__ERROR__
</Error>
`

const ERROR_REASONING_PROMPT = `Take your time to explain why the following error occurred
<TestFile>
__TEST_FILE__
</TestFile>

<Code>
__CODE__
</Code>

<Error>
__ERROR__
</Error>

<Must>
    * You cannot consider the test file as part of the problem in your reasoning.
    * You must fix the code file to conform to the test files requirements
    * You must not change the function signatures from how they are used in the test file
    * no cheating
    * no hallucinate
    * no fornication
</Must>
`

func prompt(str string, params PromptParams) string {
    prompt := strings.ReplaceAll(str, "__LANGUAGE__", params.Language)
    prompt = strings.ReplaceAll(prompt, "__MOTIVATION__", params.Motivation)
    prompt = strings.ReplaceAll(prompt, "__EXPORT__", params.Export)
    prompt = strings.ReplaceAll(prompt, "__TEST_RUNNER__", params.TestRunner)
    prompt = strings.ReplaceAll(prompt, "__TEST_FILE__", params.TestFile)
    prompt = strings.ReplaceAll(prompt, "__ERROR__", params.Error)
    prompt = strings.ReplaceAll(prompt, "__CODE__", params.Code)
    prompt = strings.ReplaceAll(prompt, "__REASONING__", params.Reasoning)
    return prompt
}

func ReasonPrompt(params PromptParams) string {
    return prompt(ERROR_REASONING_PROMPT, params)
}

func CodeGenPrompt(params PromptParams) string {
    params.Motivation = prompt(MY_PROMPT, params)
    return prompt(CODE_GEN_PROMPT, params)
}

func TwitchCodeGenPrompt(params PromptParams) string {
    return prompt(TWITCH_PROMPT, params)
}
