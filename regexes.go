// Package multiregex helps manage sets of regex rules to check against data and streams
package multiregex

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// RuleSet A set of regex rules
type RuleSet []*regexp.Regexp

// Public regex sets
var (
	Email = regexp.MustCompile(`[A-Za-z0-9_.]+((\ ?(\[|\()?\ ?@\ ?(\)|\])?\ ?)|(\ ?(\[|\()\ ?[aA][tT]\ ?(\)|\])\ ?))[0-9a-z.-]+`)
	// Bitcoin
	BitcoinAddress = regexp.MustCompile(`(?:^|[ '":])((bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39})`)
	// Hashes
	MD5    = regexp.MustCompile("\\b[A-Fa-f0-9]{32}\\b")
	Sha1   = regexp.MustCompile("\\b[A-Fa-f0-9]{40}\\b")
	Sha256 = regexp.MustCompile("\\b[A-Fa-f0-9]{64}\\b")
	Sha512 = regexp.MustCompile("\\b[A-Fa-f0-9]{128}\\b")
	// Domains
	Domain = regexp.MustCompile(`([A-Za-z0-9-]+([\[\(]?\.[\]\)]?[A-Za-z0-9-]+)*[\[\(]?\.[\]\)]?(abogado|ac|academy|accountants|active|actor|ad|adult|ae|aero|af|ag|agency|ai|airforce|al|allfinanz|alsace|am|amsterdam|an|android|ao|aq|aquarelle|ar|archi|army|arpa|as|asia|associates|at|attorney|au|auction|audio|autos|aw|ax|axa|az|ba|band|bank|bar|barclaycard|barclays|bargains|bayern|bb|bd|be|beer|berlin|best|bf|bg|bh|bi|bid|bike|bingo|bio|biz|bj|black|blackfriday|bloomberg|blue|bm|bmw|bn|bnpparibas|bo|boo|boutique|br|brussels|bs|bt|budapest|build|builders|business|buzz|bv|bw|by|bz|bzh|ca|cal|camera|camp|cancerresearch|canon|capetown|capital|caravan|cards|care|career|careers|cartier|casa|cash|cat|catering|cc|cd|center|ceo|cern|cf|cg|ch|channel|chat|cheap|christmas|chrome|church|ci|citic|city|ck|cl|claims|cleaning|click|clinic|clothing|club|cm|cn|co|coach|codes|coffee|college|cologne|com|community|company|computer|condos|construction|consulting|contractors|cooking|cool|coop|country|cr|credit|creditcard|cricket|crs|cruises|cu|cuisinella|cv|cw|cx|cy|cymru|cz|dabur|dad|dance|dating|day|dclk|de|deals|degree|delivery|democrat|dental|dentist|desi|design|dev|diamonds|diet|digital|direct|directory|discount|dj|dk|dm|dnp|do|docs|domains|doosan|durban|dvag|dz|eat|ec|edu|education|ee|eg|email|emerck|energy|engineer|engineering|enterprises|equipment|er|es|esq|estate|et|eu|eurovision|eus|events|everbank|exchange|expert|exposed|fail|farm|fashion|feedback|fi|finance|financial|firmdale|fish|fishing|fit|fitness|fj|fk|flights|florist|flowers|flsmidth|fly|fm|fo|foo|forsale|foundation|fr|frl|frogans|fund|furniture|futbol|ga|gal|gallery|garden|gb|gbiz|gd|ge|gent|gf|gg|ggee|gh|gi|gift|gifts|gives|gl|glass|gle|global|globo|gm|gmail|gmo|gmx|gn|goog|google|gop|gov|gp|gq|gr|graphics|gratis|green|gripe|gs|gt|gu|guide|guitars|guru|gw|gy|hamburg|hangout|haus|healthcare|help|here|hermes|hiphop|hiv|hk|hm|hn|holdings|holiday|homes|horse|host|hosting|house|how|hr|ht|hu|ibm|id|ie|ifm|il|im|immo|immobilien|in|industries|info|ing|ink|institute|insure|int|international|investments|io|iq|ir|irish|is|it|iwc|jcb|je|jetzt|jm|jo|jobs|joburg|jp|juegos|kaufen|kddi|ke|kg|kh|ki|kim|kitchen|kiwi|km|kn|koeln|kp|kr|krd|kred|kw|ky|kyoto|kz|la|lacaixa|land|lat|latrobe|lawyer|lb|lc|lds|lease|legal|lgbt|li|lidl|life|lighting|limited|limo|link|lk|loans|london|lotte|lotto|lr|ls|lt|ltda|lu|luxe|luxury|lv|ly|ma|madrid|maison|management|mango|market|marketing|marriott|mc|md|me|media|meet|melbourne|meme|memorial|menu|mg|mh|miami|mil|mini|mk|ml|mm|mn|mo|mobi|moda|moe|monash|money|mormon|mortgage|moscow|motorcycles|mov|mp|mq|mr|ms|mt|mu|museum|mv|mw|mx|my|mz|na|nagoya|name|navy|nc|ne|net|network|neustar|new|nexus|nf|ng|ngo|nhk|ni|ninja|nl|no|np|nr|nra|nrw|ntt|nu|nyc|nz|okinawa|om|one|ong|onl|ooo|org|organic|osaka|otsuka|ovh|pa|paris|partners|parts|party|pe|pf|pg|ph|pharmacy|photo|photography|photos|physio|pics|pictures|pink|pizza|pk|pl|place|plumbing|pm|pn|pohl|poker|porn|post|pr|praxi|press|pro|prod|productions|prof|properties|property|ps|pt|pub|pw|qa|qpon|quebec|re|realtor|recipes|red|rehab|reise|reisen|reit|ren|rentals|repair|report|republican|rest|restaurant|reviews|rich|rio|rip|ro|rocks|rodeo|rs|rsvp|ru|ruhr|rw|ryukyu|sa|saarland|sale|samsung|sarl|sb|sc|sca|scb|schmidt|schule|schwarz|science|scot|sd|se|services|sew|sexy|sg|sh|shiksha|shoes|shriram|si|singles|sj|sk|sky|sl|sm|sn|so|social|software|sohu|solar|solutions|soy|space|spiegel|sr|st|style|su|supplies|supply|support|surf|surgery|suzuki|sv|sx|sy|sydney|systems|sz|taipei|tatar|tattoo|tax|tc|td|technology|tel|temasek|tennis|tf|tg|th|tienda|tips|tires|tirol|tj|tk|tl|tm|tn|to|today|tokyo|tools|top|toshiba|town|toys|tp|tr|trade|training|travel|trust|tt|tui|tv|tw|tz|ua|ug|uk|university|uno|uol|us|uy|uz|va|vacations|vc|ve|vegas|ventures|versicherung|vet|vg|vi|viajes|video|villas|vision|vlaanderen|vn|vodka|vote|voting|voto|voyage|vu|wales|wang|watch|webcam|website|wed|wedding|wf|whoswho|wien|wiki|williamhill|wme|work|works|world|ws|wtc|wtf|xyz|yachts|yandex|ye|yoga|yokohama|youtube|yt|za|zm|zone|zuerich|zw|onion)\b)`)
	// IPs
	IPv4 = regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)([\[\(]?\.[\]\)]?)){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	IPv6 = regexp.MustCompile(`(?:[a-f0-9]{1,4}:|:){2,7}(?:[a-f0-9]{1,4}|:)`)
	// URLs
	URL = regexp.MustCompile(`(\b((http|https|hxxp|hxxps|nntp|ntp|rdp|sftp|smtp|ssh|tor|webdav|xmpp)\:\/\/[\S]+)\b)`)
	// Files
	File = regexp.MustCompile(`(([\w\-]+)\.)+(docx|doc|csv|pdf|xlsx|xls|rtf|txt|pptx|ppt|pages|keynote|numbers|exe|dll|jar|flv|swf|jpeg|jpg|gif|png|tiff|bmp|plist|app|pkg|html|htm|php|jsp|asp|zip|zipx|7z|rar|tar|gz)`)
	// Utility
	CVE = regexp.MustCompile(`(CVE-\d{4}-\d{4,7})`)

	// DefaultSet
	DefaultRules = RuleSet{Email, Domain, URL}
	MatchAll     = RuleSet{regexp.MustCompile(`.*`)}
)

// -- Functions on RuleSet --

// GetMatchedRules Given bytes return all regexes that match
func (rules RuleSet) GetMatchedRules(data []byte) RuleSet {
	matched := []*regexp.Regexp{}
	for _, rule := range rules {
		if rule.Match(data) {
			matched = append(matched, rule)
		}
	}

	return matched
}

// GetMatchedRulesReader Given a reader, return any rule matches in the stream.  Will read ENTIRE READER
// Spawns multiple go routines to check each rule
// Use limit reader to prevent reading forever
func (rules RuleSet) GetMatchedRulesReader(ctx context.Context, reader io.ReadCloser) RuleSet {
	matchedRules := make(chan *regexp.Regexp)
	finishedWorkers := make(chan bool)

	// We need to duplicate the reader stream for each worker
	// Create reader and writer for each worker thread
	workerWriters := []*io.PipeWriter{}
	workerReaders := []*io.PipeReader{}
	for range rules {
		r, w := io.Pipe()
		workerWriters = append(workerWriters, w)
		workerReaders = append(workerReaders, r)
	}

	ctxWorkers, cancelWorkers := context.WithCancel(ctx)

	// Read stream duplicator to write to all workerWriters
	go func() {
		// Defer closing all workers
		closeAll := func() {
			for _, writer := range workerWriters {
				writer.Close()
			}
		}
		defer closeAll()

		for {
			// Read
			buf := make([]byte, 1024)
			if _, err := reader.Read(buf); err != nil {
				return
			}

			// Writer to each worker stream
			for _, writer := range workerWriters {
				writer.Write(buf)
			}

			// Check if we are canceled
			select {
			case <-ctxWorkers.Done():
				return
			default:
			}
		}
	}()

	// Worker function
	workerFunction := func(workerRule *regexp.Regexp, workerReader *io.PipeReader) {
		if workerRule.MatchReader(bufio.NewReader(workerReader)) {
			// Mark this rule as matched
			select {
			case matchedRules <- workerRule:
			case <-ctxWorkers.Done():
				return
			}
		}
		// Mark us done
		select {
		case finishedWorkers <- true:
		case <-ctxWorkers.Done():
			return
		}
	}

	// Spawn worker for each rule
	for i, rule := range rules {
		go workerFunction(rule, workerReaders[i])
	}

	// Routine to capture all matches
	matches := RuleSet{}
	finishedCapturingMatches := make(chan int)
	go func() {
		for match := range matchedRules {
			matches = append(matches, match)
		}
		finishedCapturingMatches <- 0
	}()

	// Wait for threads to finish
	finishedWorkersCount := 0
	for finishedWorkersCount != len(rules) {
		select {
		case <-ctx.Done():
			cancelWorkers()
			return nil
		case <-finishedWorkers: // a go routine finished
			finishedWorkersCount++
		}
	}

	close(matchedRules)
	cancelWorkers()

	// Wait for match capture thread to finish
	<-finishedCapturingMatches

	// Return found matches
	return matches
}
