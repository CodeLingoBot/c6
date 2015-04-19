package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestLexerNext(t *testing.T) {
	l := NewLexerWithString(`.test {  }`)
	assert.NotNil(t, l)

	var r rune
	r = l.next()
	assert.Equal(t, '.', r)

	r = l.next()
	assert.Equal(t, 't', r)

	r = l.next()
	assert.Equal(t, 'e', r)

	r = l.next()
	assert.Equal(t, 's', r)

	r = l.next()
	assert.Equal(t, 't', r)
}

func TestLexerMatch(t *testing.T) {
	l := NewLexerWithString(`.foo {  }`)
	assert.NotNil(t, l)
	assert.False(t, l.match(".bar"))
	assert.True(t, l.match(".foo"))
}

func TestLexerAccept(t *testing.T) {
	l := NewLexerWithString(`.foo {  }`)
	assert.NotNil(t, l)
	assert.True(t, l.accept("."))
	assert.True(t, l.accept("f"))
	assert.True(t, l.accept("o"))
	assert.True(t, l.accept("o"))
	assert.True(t, l.accept(" "))
	assert.True(t, l.accept("{"))
}

func TestLexerIgnoreSpace(t *testing.T) {
	l := NewLexerWithString(`       .test {  }`)
	assert.NotNil(t, l)

	l.ignoreSpaces()

	var r rune
	r = l.next()
	assert.Equal(t, '.', r)

	l.backup()
	assert.True(t, l.match(".test"))
}

func TestLexerString(t *testing.T) {
	l := NewLexerWithString(`   "foo"`)
	output := l.getOutput()
	assert.NotNil(t, l)
	l.til("\"")
	lexString(l)
	token := <-output
	assert.Equal(t, T_QQ_STRING, token.Type)
}

func TestLexerTil(t *testing.T) {
	l := NewLexerWithString(`"foo"`)
	assert.NotNil(t, l)
	l.til("\"")
	assert.Equal(t, 0, l.Offset)
	l.next() // skip the quote

	l.til("\"")
	assert.Equal(t, 4, l.Offset)
}

func TestLexerAtRule(t *testing.T) {
	l := NewLexerWithString(`@import "test.css";`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_IMPORT, T_QQ_STRING, T_SEMICOLON})
	l.close()
}

func TestLexerClassNameSelector(t *testing.T) {
	l := NewLexerWithString(`.class { }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithOneProperty(t *testing.T) {
	l := NewLexerWithString(`.test { color: #fff; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTwoProperty(t *testing.T) {
	l := NewLexerWithString(`.test { color: #fff; background: #fff; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_PROPERTY_NAME, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameSelector(t *testing.T) {
	l := NewLexerWithString(`a {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameSelectorForDiv(t *testing.T) {
	l := NewLexerWithString(`div {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithUniversalSelector(t *testing.T) {
	l := NewLexerWithString(`* {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_UNIVERSAL_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`[href] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_ATTRIBUTE_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToUnquoteString(t *testing.T) {
	l := NewLexerWithString(`[lang=en] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_EQUAL, T_UNQUOTE_STRING, T_ATTRIBUTE_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToQQString(t *testing.T) {
	l := NewLexerWithString(`[lang="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_ATTRIBUTE_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorContainsQQString(t *testing.T) {
	l := NewLexerWithString(`[lang~="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_CONTAINS, T_QQ_STRING, T_ATTRIBUTE_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorAfterTagNameContainsQQString2(t *testing.T) {
	l := NewLexerWithString(`a[rel~="copyright"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_CONTAINS, T_QQ_STRING, T_ATTRIBUTE_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithMultipleAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`span[hello="Cleveland"][goodbye="Columbus"] { color: blue; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_TAGNAME_SELECTOR,
		T_AND_SELECTOR,
		T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_ATTRIBUTE_END,
		T_ATTRIBUTE_START, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_ATTRIBUTE_END,
		T_BRACE_START,
		T_PROPERTY_NAME,
		T_COLON,
		T_CONSTANT,
		T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameAndClassSelector(t *testing.T) {
	l := NewLexerWithString(`a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleUniversalSelectorPlusClassSelectorPlusAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts[href="http://google.com"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_UNIVERSAL_SELECTOR,
		T_AND_SELECTOR,
		T_CLASS_SELECTOR,
		T_AND_SELECTOR,
		T_ATTRIBUTE_START,
		T_ATTRIBUTE_NAME,
		T_EQUAL,
		T_QQ_STRING,
		T_ATTRIBUTE_END,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleUniversalPlusClassSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_UNIVERSAL_SELECTOR,
		T_AND_SELECTOR,
		T_CLASS_SELECTOR,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleChildSelector(t *testing.T) {
	l := NewLexerWithString(`div.posts > a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR,
		T_CHILD_SELECTOR,
		T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithPseudoSelector(t *testing.T) {
	var testCases = []string{`:hover {  }`, `:link {  }`, `:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []TokenType{T_PSEUDO_SELECTOR, T_BRACE_START, T_BRACE_END})
		l.close()
	}
}

func TestLexerRuleWithTagNameAndPseudoSelector(t *testing.T) {
	var testCases = []string{`a:hover {  }`, `a:link {  }`, `a:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_PSEUDO_SELECTOR, T_BRACE_START, T_BRACE_END})
		l.close()
	}
}

func TestLexerRuleLangPseudoSelector(t *testing.T) {
	// html:lang(fr-ca) { quotes: '« ' ' »' }
	l := NewLexerWithString(`html:lang(fr-ca) {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_PSEUDO_SELECTOR, T_LANG_CODE, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithIdSelector(t *testing.T) {
	l := NewLexerWithString(`#myPost {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithIdSelectorWithDigits(t *testing.T) {
	l := NewLexerWithString(`#foo123 {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithMultipleSelector(t *testing.T) {
	l := NewLexerWithString(`#foo123, .foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_COMMA, T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithVendorPrefixPropertyName(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: none; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_CONSTANT, T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithVariableAsPropertyValue(t *testing.T) {
	l := NewLexerWithString(`.test { color: $favorite; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_VARIABLE, T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerVariableAssignment(t *testing.T) {
	l := NewLexerWithString(`$favorite: #fff;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_COLON, T_HEX_COLOR, T_SEMICOLON})
	l.close()
}

func TestLexerVariableWithPtValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 10pt;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_VARIABLE, T_COLON, T_INTEGER, T_UNIT_PT, T_SEMICOLON,
	})
	l.close()
}

func TestLexerVariableWithPxValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 10px;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_VARIABLE, T_COLON, T_INTEGER, T_UNIT_PX, T_SEMICOLON,
	})
	l.close()
}

func TestLexerVariableWithEmValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 0.3em;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_VARIABLE, T_COLON, T_FLOAT, T_UNIT_EM, T_SEMICOLON,
	})
	l.close()
}

func TestLexerMultipleVariableAssignment(t *testing.T) {
	l := NewLexerWithString(`$favorite: #fff; $foo: 10em;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_VARIABLE, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_VARIABLE, T_COLON, T_INTEGER, T_UNIT_EM, T_SEMICOLON,
	})
	l.close()
}

func TestLexerSelectorWithExpansion(t *testing.T) {
	l := NewLexerWithString(`#myPost#{ abc } {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_INTERPOLATION_START, T_INTERPOLATION_END, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithSubRule(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: none;   .foo { color: #fff; } }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_CONSTANT, T_SEMICOLON,
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_BRACE_END,
		T_BRACE_END})
	l.close()
}
