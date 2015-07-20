package corpus

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/language/en"
	"github.com/blevesearch/bleve/analysis/language/it"
	"github.com/blevesearch/bleve/analysis/token_filters/lower_case_filter"
	"github.com/blevesearch/bleve/analysis/tokenizers/unicode"
	"github.com/blevesearch/bleve/registry"
)

func init() {
	registry.RegisterAnalyzer("my_base", baseAnalyzerConstructor)
	registry.RegisterAnalyzer("my_it", itAnalyzerConstructor)
	registry.RegisterAnalyzer("my_en", enAnalyzerConstructor)
}

func baseAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := bleve.Config.Cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}
	toLowerFilter, err := bleve.Config.Cache.TokenFilterNamed(lower_case_filter.Name)
	if err != nil {
		return nil, err
	}

	an := analysis.Analyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			toLowerFilter,
		},
	}
	return &an, nil
}

func itAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := bleve.Config.Cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}
	elisionFilter, err := bleve.Config.Cache.TokenFilterNamed(it.ElisionName)
	if err != nil {
		return nil, err
	}
	toLowerFilter, err := bleve.Config.Cache.TokenFilterNamed(lower_case_filter.Name)
	if err != nil {
		return nil, err
	}
	stopItFilter, err := bleve.Config.Cache.TokenFilterNamed(it.StopName)
	if err != nil {
		return nil, err
	}
	stemmerItFilter, err := bleve.Config.Cache.TokenFilterNamed(it.StemmerName)
	if err != nil {
		return nil, err
	}
	an := analysis.Analyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			elisionFilter,
			toLowerFilter,
			stopItFilter,
			stemmerItFilter,
		},
	}
	return &an, nil
}

func enAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := bleve.Config.Cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}
	possEnFilter, err := bleve.Config.Cache.TokenFilterNamed(en.PossessiveName)
	if err != nil {
		return nil, err
	}
	toLowerFilter, err := bleve.Config.Cache.TokenFilterNamed(lower_case_filter.Name)
	if err != nil {
		return nil, err
	}
	stopEnFilter, err := bleve.Config.Cache.TokenFilterNamed(en.StopName)
	if err != nil {
		return nil, err
	}
	stemmerEnFilter, err := bleve.Config.Cache.TokenFilterNamed(en.StemmerName)
	if err != nil {
		return nil, err
	}
	an := analysis.Analyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			possEnFilter,
			toLowerFilter,
			stopEnFilter,
			stemmerEnFilter,
		},
	}
	return &an, nil
}
