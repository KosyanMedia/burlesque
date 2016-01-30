PACKAGE

package cabinet
import "bitbucket.org/ww/cabinet"

Kyoto Cabinet bindings for Go. Copyleft by William Waites in 2011

 This program is free software: you can redistribute it and/or modify it under the terms of
 the GNU General Public License as published by the Free Software Foundation, either version
 3 of the License, or any later version.

 This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
 without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 See the GNU General Public License for more details.

 You should have received a copy of the GNU General Public License along with this program.
 If not, see <http://www.gnu.org/licenses/>.

Source code: http://bitbucket.org/ww/cabinet/src

Documentation: http://godoc.styx.org/pkg/bitbucket.org/ww/cabinet

These bindings have been tested with Kyoto Cabinet version 1.2.76.
They are known not to work with 1.2.7 because of the absent
kcdbreplace() function call. Once Kyoto Cabinet is installed, building
the bindings should be a simple matter of running:

    go get -u -v bitbucket.org/ww/cabinet

Simple usage will be along the lines of,

    import (
        "bitbucket.org/ww/cabinet
    )

    ...

    kc := cabinet.New()
    err = kc.Open("some_db.kch", cabinet.KCOWRITER | cabinet.KCOCREATE)
    err = kc.Set([]byte("hello"), []byte("world"))
    world, err = kc.Get([]byte("hello"))
    err = kc.Close()
    kc.Del()

Obviously checking the relevant errors...

The API follows the Kyoto Cabinet C API closely, for some examples see
http://fallabs.com/kyotocabinet/api/

Most input and output variables are []byte and not string. This is because
Kyoto Cabinet is not particularly concerned with strings and it is possible
to use any byte array as either key or value. An example from the test
suite to read an integer out of the database:

    var n int64
    v, err = kc.Get([]byte("numbers"))
    binary.Read(bytes.NewBuffer(v), binary.BigEndian, &n)

Some functions have been added for convenience using Go. The Keys()
Values() and Items() on the cursor object return a channel over which
their results will be sent, for example. This probably obviates the need
for implementing the visitor-callback pattern when using Kyoto Cabinet
with Go.

If you use this module please feel free to contact me, ww@styx.org
with any questions, comments or bug reports.


CONSTANTS

const KCEBROKEN int = C.KCEBROKEN

const KCEDUPREC int = C.KCEDUPREC

const KCEINVALID int = C.KCEINVALID

const KCELOGIC int = C.KCELOGIC

const KCEMISC = C.KCEMISC

const KCENOIMPL int = C.KCENOIMPL

const KCENOPERM int = C.KCENOPERM

const KCENOREC int = C.KCENOREC

const KCENOREPOS int = C.KCENOREPOS

const KCESUCCESS int = C.KCESUCCESS

const KCESYSTEM = C.KCESYSTEM

const KCMADD int = C.KCMADD

const KCMAPPEND int = C.KCMAPPEND

const KCMREPLACE int = C.KCMREPLACE

const KCMSET int = C.KCMSET

const KCOAUTOSYNC int = C.KCOAUTOSYNC

const KCOAUTOTRAN int = C.KCOAUTOTRAN

const KCOCREATE int = C.KCOCREATE

const KCONOLOCK int = C.KCONOLOCK

const KCONOREPAIR = C.KCONOREPAIR

const KCOREADER int = C.KCOREADER

const KCOTRUNCATE int = C.KCOTRUNCATE

const KCOTRYLOCK int = C.KCOTRYLOCK

const KCOWRITER int = C.KCOWRITER


FUNCTIONS

func EcodeName(ecode int) string


TYPES

type Item struct {
    Key   []byte
    Value []byte
}

type KCCUR struct {
    // contains unexported fields
}

func (kcc *KCCUR) Db() (kc *KCDB)

func (kcc *KCCUR) Del()

func (kcc *KCCUR) Ecode() int

func (kcc *KCCUR) Emsg() string

func (kcc *KCCUR) Get(advance bool) (k, v []byte, err os.Error)

func (kcc *KCCUR) GetKey(advance bool) (k []byte, err os.Error)

func (kcc *KCCUR) GetValue(advance bool) (v []byte, err os.Error)

func (kcc *KCCUR) Jump() (err os.Error)

func (kcc *KCCUR) JumpBack() (err os.Error)

func (kcc *KCCUR) JumpBackKey(key []byte) (err os.Error)

func (kcc *KCCUR) JumpKey(key []byte) (err os.Error)

func (kcc *KCCUR) Remove() (err os.Error)

func (kcc *KCCUR) SetValue(value []byte, advance bool) (err os.Error)

func (kcc *KCCUR) Step() (err os.Error)

func (kcc *KCCUR) StepBack() (err os.Error)

type KCDB struct {
    // contains unexported fields
}

func New() *KCDB

func (kc *KCDB) Add(key, value []byte) (err os.Error)

func (kc *KCDB) Append(key, value []byte) (err os.Error)

func (kc *KCDB) BeginTran(hard bool) (err os.Error)

func (kc *KCDB) BeginTranTry(hard bool) (err os.Error)

func (kc *KCDB) Cas(key, oval, nval []byte) (err os.Error)

func (kc *KCDB) Clear() (err os.Error)

func (kc *KCDB) Close() (err os.Error)

func (kc *KCDB) Copy(filename string) (err os.Error)

func (kc *KCDB) Count() (count uint64, err os.Error)

func (kc *KCDB) Cursor() (kcc *KCCUR)

func (kc *KCDB) Del()

func (kc *KCDB) Dump(filename string) (err os.Error)

func (kc *KCDB) Ecode() int

func (kc *KCDB) EndTran(commit bool) (err os.Error)

func (kc *KCDB) Get(key []byte) (value []byte, err os.Error)

func (kc *KCDB) IncrDouble(key []byte, amount float64) (err os.Error)

func (kc *KCDB) IncrInt(key []byte, amount int64) (err os.Error)

func (kc *KCDB) Items() (out chan *Item)

func (kc *KCDB) Keys() (out chan []byte)

func (kc *KCDB) Load(filename string) (err os.Error)

func (kc *KCDB) MatchPrefix(prefix string, max int) (matches [][]byte, err os.Error)

func (kc *KCDB) MatchRegex(regex string, max int) (matches [][]byte, err os.Error)

func (kc *KCDB) Merge(sdbs []*KCDB, mode int) (err os.Error)

func (kc *KCDB) Open(filename string, mode int) (err os.Error)

func (kc *KCDB) Path() (path string, err os.Error)

func (kc *KCDB) Remove(key []byte) (err os.Error)

func (kc *KCDB) Replace(key, value []byte) (err os.Error)

func (kc *KCDB) Set(key, value []byte) (err os.Error)

func (kc *KCDB) Size() (size uint64, err os.Error)

func (kc *KCDB) Status() (status string, err os.Error)

func (kc *KCDB) Sync(hard bool) (err os.Error)

func (kc *KCDB) Values() (out chan []byte)

func Version() string

SUBDIRECTORIES

	.hg
