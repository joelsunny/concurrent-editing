import React from 'react';
import './App.css';
import Quill from 'quill';
import Delta from 'quill-delta';

var io = new WebSocket("ws://localhost:8080");
var quill;

// client side state variables
var index = 0;
var lastIndex = 0;
var prevChar = '';

class Editor extends React.Component {

    constructor(props) {
        super(props);

        this.log = this.log.bind(this);
    }

    componentDidMount() {
        var container = document.getElementById('e' + this.props.id);
        quill = new Quill(container, {
            modules:{
                toolbar:false
            },
            placeholder: 'Compose an epic...',
            theme: 'snow'  // or 'bubble'
          });
          // add a listener for text-change event
          quill.on('text-change', this.log);

    }

    async deltaPropogate(msg) {
        io.send(msg);
        index += 1
    }

    serverSync() {

    }

    log(delta, oldDelta, source) {
        const logstream = document.getElementById("c" + this.props.id);
        logstream.innerHTML = "delta : " + JSON.stringify(delta) + "<br>" + "old: " + JSON.stringify(oldDelta) + "<br>" + "source: " + source + "<br>" + "-------------------<br><br>" + logstream.innerHTML;
        if (source === "user" ) {

            let m = {
                type: "delta",
                op: {
                    retain:0,
                    insert:'',
                    delete:0,
                },
                index: index,
                lastIndex: lastIndex
            }

            for (let i = 0; i < delta["ops"].length; i++) {
                if (delta["ops"][i].retain !== undefined) {
                    m.op.retain = delta["ops"][i].retain
                }
                if (delta["ops"][i].insert !== undefined) {
                    let ins = delta["ops"][i].insert
                    m.op.insert = ins
                    // if (prevChar == " ") {
                    //     m.op.insert = " " + ins
                    // } else {
                    //     m.op.insert = ins
                    // }
                    // prevChar = ins
                }
                if (delta["ops"][i].delete !== undefined) {
                    m.op.delete = delta["ops"][i].delete
                }
            }
            this.deltaPropogate(JSON.stringify(m));
        }
        
    }

    render() {
        return(
            <div id={this.props.id} className="editor">
                <div className="controls">
                    controls
                </div>
                <div id={"e" + this.props.id} className="text-area">
                </div>
                <div id={"c" + this.props.id} className="console">
                    
                </div>
            </div>
        )
    }
}

class ServerView extends Editor {
    constructor(props) {
        super(props);
    }

    componentDidMount() {

        io.onopen = () => {
                            this.log("url: " + io.url);
                            this.log("connection status: " + io.readyState);
                        };
        
        io.onmessage = (msg) => {
                this.log(msg.data);
                let json_msg = JSON.parse(msg.data)
                console.log(json_msg)
                let MsgType = json_msg["type"]
                console.log("MsgType : " + MsgType)
                if (MsgType === undefined) {
                    console.log("undefined message type")
                } else if(MsgType === "delta") {
                    let delta = json_msg["log"]["delta"]
                    console.log("INFO: received delta from server" )
                    //console.log(JSON.stringify(delta))
                    console.log(delta.insert.length)
                    quill.updateContents(new Delta()
                                                .retain(delta["retain"])
                                                .delete(delta["delete"])
                                                .insert(delta["insert"])
                                                , 'api')
                    lastIndex = json_msg["log"]["index"]

                } else if(MsgType === "sync") {
                    if (json_msg["log"] === null) {
                        return
                    }
                    for (let i = 0; i < json_msg["log"].length; i++) {
                        let delta = json_msg["log"][i]["delta"]
                        quill.updateContents(new Delta()
                                                .retain(delta["retain"])
                                                .delete(delta["delete"])
                                                .insert(delta["insert"])
                                                , 'api')
                        lastIndex = json_msg["log"]["index"]
                    }
                } else {
                    console.log(MsgType)
                }
                
        };

        io.onclose = () => {
            this.log("connection terminated");
        };
    }

    log(msg) {
        const logstream = document.getElementById("c" + this.props.id);
        logstream.innerHTML = "got: " + msg + "<br>" + "<br>" + logstream.innerHTML;
    }

    render() {
        return(
            <div id={this.props.id} className="editor">
                <div id={"c" + this.props.id} className="server-console">
                    
                </div>
            </div>
        )
    }
}

function App() {
    return (
    <div className="App">
        <div className="logo-wrapper">
            <img className="logo" src="logo.png"/>
        </div>
        <Editor id="1" />
        <ServerView id="2" />
    </div>
    );
}

export default App;
