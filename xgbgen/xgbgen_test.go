package xgbgen

import (
	"strings"
	"testing"
)

var testinput = `
/* Definitions for the X window system used by server and c bindings */

/***********************************************************

Copyright 1987, 1998  The Open Group

Permission to use, copy, modify, distribute, and sell this software and its
documentation for any purpose is hereby granted without fee, provided that
the above copyright notice appear in all copies and that both that
copyright notice and this permission notice appear in supporting
documentation.

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
OPEN GROUP BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Except as contained in this notice, the name of The Open Group shall not be
used in advertising or otherwise to promote the sale, use or other dealings
in this Software without prior written authorization from The Open Group.


Copyright 1987 by Digital Equipment Corporation, Maynard, Massachusetts.

                        All Rights Reserved

Permission to use, copy, modify, and distribute this software and its
documentation for any purpose and without fee is hereby granted,
provided that the above copyright notice appear in all copies and that
both that copyright notice and this permission notice appear in
supporting documentation, and that the name of Digital not be
used in advertising or publicity pertaining to distribution of the
software without specific, written prior permission.

DIGITAL DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE, INCLUDING
ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO EVENT SHALL
DIGITAL BE LIABLE FOR ANY SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR
ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS,
WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION,
ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS
SOFTWARE.

******************************************************************/

#define XPROTO_H

#define KeySym CARD32

#define X_TCP_PORT 6000     /* add display number */

#define xTrue        1
#define xFalse       0


typedef CARD16 KeyButMask;

/*****************
   Connection setup structures.  See Chapter 8: Connection Setup
   of the X Window System Protocol specification for details.
*****************/

/* Client initiates handshake with this data, followed by the strings
 * for the auth protocol & data.
 */
typedef struct {
    CARD8	byteOrder;
    BYTE	pad;
    CARD16	majorVersion, minorVersion;
    CARD16	nbytesAuthProto;	/* Authorization protocol */
    CARD16	nbytesAuthString;	/* Authorization string */
    CARD16	pad2;
} xConnClientPrefix;
`

func TestLexer(t *testing.T) {

	filehandle := strings.NewReader(testinput)
	lxr := lex("testinput", filehandle)

	t.Run("test scan", func(t *testing.T) {
		for tkn := lxr.nextToken(); tkn.Type != tknEOF; tkn = lxr.nextToken() {
			if tkn.Type == tknError {
				t.Log(tkn.String())
				break
			}
			t.Log(tkn.String())
		}
	})
}
