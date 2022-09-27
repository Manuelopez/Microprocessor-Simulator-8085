import { useEffect, useState } from 'react'
import './App.css'
import CodeMirror from '@uiw/react-codemirror'
import { javascript } from '@codemirror/lang-javascript'
import { okaidia } from '@uiw/codemirror-theme-okaidia'

const socket = new WebSocket('ws://127.0.0.1:8080/ws')

const binaryToHex = (val: number[]) => {
  const s = val.join('')
  const d = parseInt(s, 2)
  const sd = d.toString(16)
  return '0x' + ('0'.repeat(Math.max(4 - sd.length, 0)) + sd)
}

const binaryToDecimal = (val: number[]) => {
  const s = val.join('')
  const d = parseInt(s, 2)
  return d.toString(10)
}
function App() {
  const [instructions, setInstructions] = useState('')
  const [lineNumbers, setLineNumbers] = useState(0)
  const [microprocessor, setMicroProcessor] = useState<{
    al: number[]
    ah: number[]
    b: number[]
    c: number[]
    d: number[]
    e: number[]
    l: number[]
    h: number[]
    ir: number[]
    pc: number[]
    mar: number[]
    mbr: number[]
    memory: number[][][]
    stack: {
      mem: number[][][]
      sp: number[]
      top: number[]
      bottom: number[]
    }
    alu: {
      temp1: number[]
      temp2: number[]
      carry: boolean
      zero: boolean
      comparison: number[]
    }
  }>()
  const [registers, setRegisters] = useState(true)

  const [formatAl, setFormatAl] = useState('b')
  const [formatAh, setFormatAh] = useState('b')
  const [formatB, setFormatB] = useState('b')
  const [formatC, setFormatC] = useState('b')
  const [formatD, setFormatD] = useState('b')
  const [formatE, setFormatE] = useState('b')
  const [formatL, setFormatL] = useState('b')
  const [formatH, setFormatH] = useState('b')
  const [formatIr, setFormatIr] = useState('b')
  const [formatPc, setFormatPc] = useState('b')
  const [formatMar, setFormatMar] = useState('b')
  const [formatMbr, setFormatMbr] = useState('b')
  const [formatTop, setFormatTop] = useState('b')
  const [formatBottom, setFormatBottom] = useState('b')
  const [formatSp, setFormatSp] = useState('b')
  const [formatTemp1, setFormatTemp1] = useState('b')
  const [formatTemp2, setFormatTemp2] = useState('b')
  const [formatComparison, setFormatComparison] = useState('b')

  let a = document.getElementById('sp')
  a?.scrollIntoView()

  useEffect(() => {
    socket.onopen = () => {
      console.log('Successfully connect')
      socket.send(JSON.stringify({ type: 'connecting' }))
    }
    socket.onclose = (event) => {
      socket.send(JSON.stringify({ type: 'disconecting' }))
    }
    socket.onerror = (err) => {
      console.log('Socket Error: ', err)
    }
    socket.onmessage = (message) => {
      let data = JSON.parse(message.data)
      setMicroProcessor(data)
    }
  })

  return (
    <div className='App flex flex-row'>
      <div>
        <div></div>
      </div>

      <div className='h-screen'>
        <button
          onClick={() => {
            socket.send(
              JSON.stringify({
                type: 'start',
                data: { instructions: instructions.split('\n') },
              })
            )
          }}
        >
          Send
        </button>

        <CodeMirror
          value={instructions}
          height='calc(100vh - 50px)'
          width='30vw'
          theme={okaidia}
          extensions={[javascript({ typescript: true })]}
          basicSetup={{
            lineNumbers: true,
            lintKeymap: false,
          }}
          onChange={(text) => {
            setInstructions(text)
          }}
        />
      </div>
      <div className='flex flex-col'>
        <div className=' flex flex-row h-60 '>
          <div>
            <ul className=' flex flex-row bg-gray-300 h-fit w-fit border-2 border-black'>
              <li
                className='w-20 ml-2 h-10 leading-10 border-r-2 border-black'
                onClick={() => setRegisters(true)}
              >
                Registers
              </li>
              <li
                className={`mr-2 ml-3  h-10 leading-10`}
                onClick={() => setRegisters(false)}
              >
                Flags
              </li>
            </ul>
            {registers ? (
              <>
                <ul className=' rounded  w-96 flex flex-row flex-wrap '>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatAl === 'b') {
                        setFormatAl('h')
                      } else if (formatAl === 'h') {
                        setFormatAl('d')
                      } else {
                        setFormatAl('b')
                      }
                    }}
                  >
                    AL <br />
                    {microprocessor
                      ? formatAl === 'b'
                        ? microprocessor.al
                        : formatAl === 'h'
                        ? binaryToHex(microprocessor.al)
                        : binaryToDecimal(microprocessor.al)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatAh === 'b') {
                        setFormatAh('h')
                      } else if (formatAh === 'h') {
                        setFormatAh('d')
                      } else {
                        setFormatAh('b')
                      }
                    }}
                  >
                    AH <br />
                    {microprocessor
                      ? formatAh === 'b'
                        ? microprocessor.ah
                        : formatAh === 'h'
                        ? binaryToHex(microprocessor.ah)
                        : binaryToDecimal(microprocessor.ah)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatB === 'b') {
                        setFormatB('h')
                      } else if (formatB === 'h') {
                        setFormatB('d')
                      } else {
                        setFormatB('b')
                      }
                    }}
                  >
                    B <br />
                    {microprocessor
                      ? formatB === 'b'
                        ? microprocessor.b
                        : formatB === 'h'
                        ? binaryToHex(microprocessor.b)
                        : binaryToDecimal(microprocessor.b)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatC === 'b') {
                        setFormatC('h')
                      } else if (formatC === 'h') {
                        setFormatC('d')
                      } else {
                        setFormatC('b')
                      }
                    }}
                  >
                    C <br />
                    {microprocessor
                      ? formatC === 'b'
                        ? microprocessor.c
                        : formatC === 'h'
                        ? binaryToHex(microprocessor.c)
                        : binaryToDecimal(microprocessor.c)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatD === 'b') {
                        setFormatD('h')
                      } else if (formatD === 'h') {
                        setFormatD('d')
                      } else {
                        setFormatD('b')
                      }
                    }}
                  >
                    D <br />
                    {microprocessor
                      ? formatD === 'b'
                        ? microprocessor.d
                        : formatD === 'h'
                        ? binaryToHex(microprocessor.d)
                        : binaryToDecimal(microprocessor.d)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatE === 'b') {
                        setFormatE('h')
                      } else if (formatE === 'h') {
                        setFormatE('d')
                      } else {
                        setFormatE('b')
                      }
                    }}
                  >
                    E <br />
                    {microprocessor
                      ? formatE === 'b'
                        ? microprocessor.e
                        : formatE === 'h'
                        ? binaryToHex(microprocessor.e)
                        : binaryToDecimal(microprocessor.e)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatH === 'b') {
                        setFormatH('h')
                      } else if (formatH === 'h') {
                        setFormatH('d')
                      } else {
                        setFormatH('b')
                      }
                    }}
                  >
                    H <br />
                    {microprocessor
                      ? formatH === 'b'
                        ? microprocessor.h
                        : formatH === 'h'
                        ? binaryToHex(microprocessor.h)
                        : binaryToDecimal(microprocessor.h)
                      : 'a'}
                  </li>
                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatL === 'b') {
                        setFormatL('h')
                      } else if (formatL === 'h') {
                        setFormatL('d')
                      } else {
                        setFormatL('b')
                      }
                    }}
                  >
                    L <br />
                    {microprocessor
                      ? formatL === 'b'
                        ? microprocessor.l
                        : formatL === 'h'
                        ? binaryToHex(microprocessor.l)
                        : binaryToDecimal(microprocessor.l)
                      : 'a'}
                  </li>

                  <li
                    className='p-1'
                    onClick={() => {
                      if (formatIr === 'b') {
                        setFormatIr('h')
                      } else if (formatIr === 'h') {
                        setFormatIr('d')
                      } else {
                        setFormatIr('b')
                      }
                    }}
                  >
                    IR <br />
                    {microprocessor
                      ? formatIr === 'b'
                        ? microprocessor.ir
                        : formatIr === 'h'
                        ? binaryToHex(microprocessor.ir)
                        : binaryToDecimal(microprocessor.ir)
                      : 'a'}
                  </li>
                </ul>
              </>
            ) : (
              <ul className='rounded  w-96 flex flex-row flex-wrap '>
                <li
                  className='p-1 pr-1'
                  onClick={() => {
                    if (formatTemp1 === 'b') {
                      setFormatTemp1('h')
                    } else if (formatTemp1 === 'h') {
                      setFormatTemp1('d')
                    } else {
                      setFormatTemp1('b')
                    }
                  }}
                >
                  Temp1 <br />
                  {microprocessor
                    ? formatTemp1 === 'b'
                      ? microprocessor.alu.temp1
                      : formatTemp1 === 'h'
                      ? binaryToHex(microprocessor.alu.temp1)
                      : binaryToDecimal(microprocessor.alu.temp1)
                    : 'a'}
                </li>
                <li
                  className='p-1 pr-44'
                  onClick={() => {
                    if (formatTemp2 === 'b') {
                      setFormatTemp2('h')
                    } else if (formatTemp2 === 'h') {
                      setFormatTemp2('d')
                    } else {
                      setFormatTemp2('b')
                    }
                  }}
                >
                  Temp2 <br />
                  {microprocessor
                    ? formatTemp2 === 'b'
                      ? microprocessor.alu.temp2
                      : formatTemp2 === 'h'
                      ? binaryToHex(microprocessor.alu.temp2)
                      : binaryToDecimal(microprocessor.alu.temp2)
                    : 'a'}
                </li>
                <li className='p-1 pr-1'>
                  Carry <br /> {microprocessor?.alu.carry ? 'True' : 'False'}
                </li>
                <li className='p-1 pr-48'>
                  Zero <br /> {microprocessor?.alu.zero ? 'True' : 'False'}
                </li>
                <li
                  className='p-1 pr-8'
                  onClick={() => {
                    if (formatComparison === 'b') {
                      setFormatComparison('h')
                    } else if (formatComparison === 'h') {
                      setFormatComparison('d')
                    } else {
                      setFormatComparison('b')
                    }
                  }}
                >
                  Comparison <br />
                  {microprocessor
                    ? formatComparison === 'b'
                      ? microprocessor.alu.comparison
                      : formatComparison === 'h'
                      ? binaryToHex(microprocessor.alu.comparison)
                      : binaryToDecimal(microprocessor.alu.comparison)
                    : 'a'}
                </li>
              </ul>
            )}
          </div>
          <div>
            <h1>STACK</h1>

            <div className='overflow-y-scroll snap-end h-48'>
              {microprocessor?.stack.mem
                .map((segment, index1) => {
                  let sp = microprocessor?.stack.sp
                  let hbit = 0
                  let lbit = 0
                  if (sp != null) {
                    let val = parseInt(binaryToDecimal(sp))
                    console.log(val)
                    hbit = val >> 3
                    lbit = val & 0x7
                  }
                  return segment
                    .map((number, index2) => {
                      if (sp != null && index1 == hbit && index2 == lbit) {
                        return (
                          <div className='flex flex-row'>
                            <div
                              className='snap-end'
                              key={`${index1} - ${index2}`}
                            >
                              {number}
                            </div>
                            <div id='sp' key={`SP ${index1} - ${index2}`}>
                              {'<---- SP'}
                            </div>
                          </div>
                        )
                      }
                      return (
                        <div>
                          <div key={`${index1} - ${index2}`}>{number}</div>
                        </div>
                      )
                    })
                    .slice(0)
                    .reverse()
                })
                .slice(0)
                .reverse()}
            </div>
          </div>
          <div className='p-6'>
            <h1
              onClick={() => {
                if (formatTop === 'b') {
                  setFormatTop('h')
                } else if (formatTop === 'h') {
                  setFormatTop('d')
                } else {
                  setFormatTop('b')
                }
              }}
            >
              Top <br />
              {microprocessor
                ? formatTop === 'b'
                  ? microprocessor.stack.top
                  : formatTop === 'h'
                  ? binaryToHex(microprocessor.stack.top)
                  : binaryToDecimal(microprocessor.stack.top)
                : 'a'}
            </h1>
            <h1
              onClick={() => {
                if (formatBottom === 'b') {
                  setFormatBottom('h')
                } else if (formatBottom === 'h') {
                  setFormatBottom('d')
                } else {
                  setFormatBottom('b')
                }
              }}
            >
              Bottom <br />
              {microprocessor
                ? formatBottom === 'b'
                  ? microprocessor.stack.bottom
                  : formatBottom === 'h'
                  ? binaryToHex(microprocessor.stack.bottom)
                  : binaryToDecimal(microprocessor.stack.bottom)
                : 'a'}
            </h1>
          </div>
        </div>
        <div>
          <div className='flex flex-col h-96 w-96 overflow-x-scroll  flex-wrap'>
            {microprocessor?.memory.map((segment) => {
              return (
                <div className='p-1'>
                  {segment.map((number) => (
                    <div>{binaryToHex(number)}</div>
                  ))}
                </div>
              )
            })}
          </div>
          <ul>
            <li
              onClick={() => {
                if (formatMar === 'b') {
                  setFormatMar('h')
                } else if (formatMar === 'h') {
                  setFormatMar('d')
                } else {
                  setFormatMar('b')
                }
              }}
            >
              MAR <br />
              {microprocessor
                ? formatMar === 'b'
                  ? microprocessor.mar
                  : formatMar === 'h'
                  ? binaryToHex(microprocessor.mar)
                  : binaryToDecimal(microprocessor.mar)
                : 'a'}
            </li>
            <li
              onClick={() => {
                if (formatMbr === 'b') {
                  setFormatMbr('h')
                } else if (formatMbr === 'h') {
                  setFormatMbr('d')
                } else {
                  setFormatMbr('b')
                }
              }}
            >
              MBR <br />
              {microprocessor
                ? formatMbr === 'b'
                  ? microprocessor.mbr
                  : formatMbr === 'h'
                  ? binaryToHex(microprocessor.mbr)
                  : binaryToDecimal(microprocessor.mbr)
                : 'a'}
            </li>
            <li
              onClick={() => {
                if (formatPc === 'b') {
                  setFormatPc('h')
                } else if (formatPc === 'h') {
                  setFormatPc('d')
                } else {
                  setFormatPc('b')
                }
              }}
            >
              PC <br />
              {microprocessor
                ? formatPc === 'b'
                  ? microprocessor.pc
                  : formatPc === 'h'
                  ? binaryToHex(microprocessor.pc)
                  : binaryToDecimal(microprocessor.pc)
                : 'a'}
            </li>
          </ul>
        </div>
      </div>
    </div>
  )
}

export default App
