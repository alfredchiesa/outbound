package outbound

// Some basic structs for now
type Response struct {
	StatusCode int
	Status     string
	Protocol   string
}

type Request struct {
	Uri         string
	Accept      string
	Method      string
	Body        string
	UserAgent   string
	ContentType string
	Host        string
	Headers     []string
}

// Public Func Prototypes
// these will be like outbound.GET, outbound.PUT, outbound.UDP
func GET(...interface{}) {
	//pass
}

func POST(...interface{}) {
	//pass
}

func PUT(...interface{}) {
	//pass
}

func DELETE(...interface{}) {
	//pass
}

func HEAD(...interface{}) {
	//pass
}

func OPTIONS(...interface{}) {
	//pass
}

func PATCH(...interface{}) {
	//pass
}

func UDP(...interface{}) {
	//pass
}
