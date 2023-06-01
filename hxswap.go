package htmxtools

// HXSwap - https://htmx.org/attributes/hx-swap/
type HXSwap int64

const (
	// SwapUnknown is the zero value for the enum
	SwapUnknown HXSwap = iota
	// SwapInnerHTML - The default, replace the inner html of the target element - probably no reason to ever set this explicitly
	SwapInnerHTML
	// SwapOuterHTML - Replace the entire target element with the response
	SwapOuterHTML
	// SwapBeforeBegin - Insert the response before the target element
	SwapBeforeBegin
	// SwapAfterBegin - Insert the response before the first child of the target element
	SwapAfterBegin
	// SwapBeforeEnd - Insert the response after the last child of the target element
	SwapBeforeEnd
	// SwapAfterEnd - Insert the response after the target element
	SwapAfterEnd
	// SwapDelete - Deletes the target element regardless of the response
	SwapDelete
	// SwapNone - Does not append content from response (out of band items will still be processed).
	SwapNone
)

// HXSwapFromString returns an [HXSWap] from its string representation
func HXSwapFromString(s string) HXSwap {
	switch s {
	case SwapInnerHTML.String():
		return SwapInnerHTML
	case SwapOuterHTML.String():
		return SwapOuterHTML
	case SwapBeforeBegin.String():
		return SwapBeforeBegin
	case SwapAfterBegin.String():
		return SwapAfterBegin
	case SwapAfterEnd.String():
		return SwapAfterEnd
	case SwapDelete.String():
		return SwapDelete
	case SwapNone.String():
		return SwapNone
	default:
		return SwapInnerHTML
	}
}

// String returns the string representation of a status code
func (hxs HXSwap) String() string {
	switch hxs {
	case SwapInnerHTML:
		return "innerHTML"
	case SwapOuterHTML:
		return "outerHTML"
	case SwapBeforeBegin:
		return "beforebegin"
	case SwapAfterBegin:
		return "afterbegin"
	case SwapBeforeEnd:
		return "beforeend"
	case SwapAfterEnd:
		return "afterend"
	case SwapDelete:
		return "delete"
	case SwapNone:
		return "none"
	default:
		return "innerHTML"
	}
}
