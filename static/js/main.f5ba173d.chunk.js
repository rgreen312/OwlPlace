(window.webpackJsonpclient=window.webpackJsonpclient||[]).push([[0],{205:function(e,t,n){e.exports=n(439)},210:function(e,t,n){},211:function(e,t,n){},380:function(e,t,n){},424:function(e,t,n){},425:function(e,t,n){},438:function(e,t,n){},439:function(e,t,n){"use strict";n.r(t);var a,r,o=n(0),c=n.n(o),i=n(6),u=(n(210),n(91)),l=n(92),s=n(96),f=n(93),p=n(95),m=(n(211),n(58)),d=n(41),g=Object(d.f)(function(e){var t=e.children,n=e.location.pathname;return Object(o.useEffect)(function(){window.scrollTo(0,0)},[n]),t||null}),v=n(29),b=n(59),O=n(203),h=n(191),E=n(98),y=(n(380),function(e){var t=e.onComplete,n=e.onCancel,a=Object(o.useState)({r:0,g:0,b:0}),r=Object(O.a)(a,2),i=r[0],u=r[1];return c.a.createElement("div",null,c.a.createElement(h.SketchPicker,{color:i,onChange:function(e){return u({r:e.rgb.r,b:e.rgb.b,g:e.rgb.g})}}),c.a.createElement("div",{className:"button-bar"},c.a.createElement(E.a,{onClick:function(){return t(i)},className:"okay-button"},"Okay"),c.a.createElement(E.a,{onClick:n,className:"cancel-button"},"Cancel")))}),j=(n(424),n(13)),C=n(3),k=n.n(C),w=function(e){function t(e){var n;return Object(u.a)(this,t),(n=Object(s.a)(this,Object(f.a)(t).call(this,e))).canvasRef=void 0,n.canvasRef=Object(o.createRef)(),n.state={showColorPicker:!1,isDrag:!1,translateX:0,translateY:0,dragStartX:0,dragStartY:0},n.onCancel=n.onCancel.bind(Object(b.a)(n)),n.onComplete=n.onComplete.bind(Object(b.a)(n)),n.updateTranslate=n.updateTranslate.bind(Object(b.a)(n)),n}return Object(p.a)(t,e),Object(l.a)(t,[{key:"componentDidMount",value:function(){var e=this;this.canvasRef.current.width=1e3,this.canvasRef.current.height=1e3;var t=this.canvasRef.current.getContext("2d");t&&this.props.registerContext(t),t.imageSmoothingEnabled=!1,t.fillStyle="#000000",t.fillRect(0,0,1e3,500),t.fillStyle="#ff0000",t.fillRect(0,500,1e3,500),this.canvasRef.current.addEventListener("mousemove",function(t){if(!e.state.showColorPicker){var n=e.getMousePos(e.canvasRef.current,t),a=n.x,r=n.y;e.props.updatePosition(a,r)}}),this.canvasRef.current.addEventListener("mouseout",function(){e.state.showColorPicker||e.props.onMouseOut()}),this.canvasRef.current.addEventListener("mousedown",function(t){var n=e.state,a=n.translateX,r=n.translateY,o=t.clientX-a,c=t.clientY-r;e.setState({dragStartX:o,dragStartY:c}),e.canvasRef.current.addEventListener("mousemove",e.updateTranslate)}),this.canvasRef.current.addEventListener("mouseup",function(t){if(e.canvasRef.current.removeEventListener("mousemove",e.updateTranslate),!e.state.isDrag){var n=e.getMousePos(e.canvasRef.current,t),a=n.x,r=n.y;e.props.updatePosition(a,r),e.showColorPicker()}e.setState({isDrag:!1})})}},{key:"updateTranslate",value:function(e){var t=this.state,n=t.dragStartX,a=t.dragStartY,r=e.clientX-n,o=e.clientY-a;this.setState({isDrag:!0,translateX:r,translateY:o})}},{key:"getMousePos",value:function(e,t){var n=e.getBoundingClientRect();return{x:t.clientX-n.left,y:t.clientY-n.top}}},{key:"onCancel",value:function(){this.hideColorPicker()}},{key:"onComplete",value:function(e){this.hideColorPicker();var t=this.canvasRef.current.getContext("2d"),n=this.props.position.x-1,a=this.props.position.y-1;t.fillStyle="rgb(".concat(e.r,", ").concat(e.g,", ").concat(e.b,")"),t.fillRect(n,a,1,1),this.props.onUpdatePixel({r:e.r,g:e.g,b:e.b},n,a)}},{key:"showColorPicker",value:function(){this.setState({showColorPicker:!0})}},{key:"hideColorPicker",value:function(){this.setState({showColorPicker:!1})}},{key:"render",value:function(){var e=this,t=this.props,n=(t.receivedError,t.zoomFactor),a=t.setZoomFactor,r=this.state,o=r.translateX,i=r.translateY,u=r.isDrag;return c.a.createElement("div",{className:"canvas-container"},this.state.showColorPicker&&c.a.createElement("div",{className:"color-picker"},c.a.createElement(y,{onCancel:this.onCancel,onComplete:function(t){return e.onComplete(t)}})),c.a.createElement("div",{className:k()({"pan-canvas":!0,"drag-canvas":u}),style:{transform:"translate(".concat(o,"px, ").concat(i,"px)")}},c.a.createElement("div",{className:"zoom-canvas",style:{transform:"scale(".concat(n,", ").concat(n,")")}},c.a.createElement("canvas",{ref:this.canvasRef}))),c.a.createElement("div",{className:"zoom-controls"},c.a.createElement(j.a,{type:"plus-circle",onClick:function(){return a(n+10)},className:"zoom-icon"}),c.a.createElement(j.a,{type:"minus-circle",onClick:function(){return a(n-10)},className:"zoom-icon"})))}}]),t}(o.Component),S=function(e){return e.websocket.socket||void 0},N=function(e){return e.websocket.receivedError},P=function(e){return e.canvas.curPosition||void 0},R=function(e){return e.canvas.zoomFactor},T={registerContext:function(e){return function(t){t(function(e){return{type:"canvas/REGISTER_CONTEXT",payload:{context:e}}}(e))}},updatePosition:function(e,t){return function(n,a){var r=a(),o=R(r);n(function(e,t){return{type:"canvas/UPDATE_POSITION",payload:{x:e,y:t}}}(Math.ceil(e/o),Math.ceil(t/o)))}},onMouseOut:function(){return function(e){e({type:"canvas/CLEAR_POSITION"})}},setZoomFactor:function(e){return function(t,n){if(!(e<0)){var a=function(e){return e.canvas.canvasContext||void 0}(n());a&&a.scale(e,e),t({type:"canvas/SET_ZOOM",payload:{zoom:e}})}}},onUpdatePixel:function(e,t,n){return function(e){}}},x=Object(v.b)(function(e){return{receivedError:N(e),zoomFactor:R(e),position:P(e)}},T)(w),I=function(){return c.a.createElement(x,null)},L="websocket/START_CONNECT",_="websocket/CONNECT_SUCCESS",A="websocket/CONNECT_ERROR",M="websocket/CLOSE_CONNNECTION",X=function(){return{type:M}},U=function(e){var t=e.sendUpdateMessage;return c.a.createElement("div",{className:"about-page"},"Update the canvas one pixel at a time...",c.a.createElement("h2",null,"change the canvas one pixel at a time"),c.a.createElement("p",null,'Click "Pixel 1" to send an update message to the server!'),c.a.createElement("button",{onClick:function(){return t("user1",10,400,255,255,255)},id:"p1"}," Pixel 1 "))},z={sendUpdateMessage:function(e,t,n,a,r,o){return function(c,i){var u=S(i());u&&u.send(function(e,t,n,a,r,o){return JSON.stringify({type:1,userId:e,x:t,y:n,r:a,g:r,b:o})}(e,t,n,a,r,o))}}},D=Object(v.b)(null,z)(U),Y=n(138),F=n(442),G=n(440),J=n(441),B=(n(425),function(e){var t=e.onLogin,n=e.isLoggedIn,a=e.name,r=e.onLogout;window.onGoogleScriptLoad=function(){console.log("The google script has really loaded, cool!")};var o=n?c.a.createElement("div",{className:"name-label"},"Hi, ",a):c.a.createElement(E.a,{className:"login-button",onClick:t},"Login"),i=Object(d.e)().pathname,u=c.a.createElement(Y.a,null,c.a.createElement(Y.a.Item,null,"// Check the last 5 characters in a string.","about"!=i.substring(i.length-5,i.length)?c.a.createElement(m.b,{to:"/about"},"About"):c.a.createElement(m.b,{to:"/"},"Home")),n&&c.a.createElement(Y.a.Item,{onClick:r},"Sign Out")),l=c.a.createElement(F.a,{key:"more",overlay:u},c.a.createElement(E.a,{style:{border:"none",padding:0}},c.a.createElement(j.a,{type:"ellipsis",style:{fontSize:20,verticalAlign:"top"}})));return c.a.createElement(G.a,{title:"OwlPlace",subTitle:"change the canvas one pixel at a time",tags:c.a.createElement(J.a,{color:"green"},"COMP 413"),extra:[o,l]})}),Z=n(50),H=n.n(Z),q=n(97),V="login/LOGIN_START",W="login/LOGIN_SUCCESS",K="login/LOGIN_ERROR",Q="login/SIGN_OUT",$=function(){return{type:V}},ee=function(e,t,n){return{type:W,payload:{name:e,id:t,email:n}}},te=function(e){return e.login.isLoggedIn},ne=function(e){return e.login.name||void 0},ae={onLogin:function(){return function(){var e=Object(q.a)(H.a.mark(function e(t){var n,a,r;return H.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t($()),n=new Promise(function(e){gapi.load("auth2",function(){gapi.auth2.init({client_id:"634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com"}),e()})}),e.next=4,n;case 4:return e.next=6,gapi.auth2.getAuthInstance().signIn();case 6:a=gapi.auth2.getAuthInstance().currentUser.get(),r=a.getBasicProfile(),t(ee(r.getName(),r.getId(),r.getEmail()));case 9:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}()},onLogout:function(){return function(){var e=Object(q.a)(H.a.mark(function e(t){var n;return H.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return n=gapi.auth2.getAuthInstance(),e.next=3,n.signOut();case 3:t({type:Q});case 4:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}()}},re=Object(v.b)(function(e){return{isLoggedIn:te(e),name:ne(e)}},ae)(B),oe=(n(438),function(){return c.a.createElement("div",null,c.a.createElement("footer",null,"Made with \u2665 by COMP 413 @ Rice"))}),ce=function(){return c.a.createElement("div",{className:"error-page"},"There was an error.")},ie=function(){return c.a.createElement(m.a,{basename:"/OwlPlace"},c.a.createElement(g,null,c.a.createElement(re,null),c.a.createElement(d.a,{exact:!0,path:"/",component:I}),c.a.createElement(d.a,{path:"/about",component:D}),c.a.createElement(d.a,{path:"/error",component:ce}),c.a.createElement(oe,null)))},ue=function(e){function t(){return Object(u.a)(this,t),Object(s.a)(this,Object(f.a)(t).apply(this,arguments))}return Object(p.a)(t,e),Object(l.a)(t,[{key:"componentDidMount",value:function(){this.props.checkLogin(),this.props.openConnection()}},{key:"render",value:function(){return c.a.createElement(ie,null)}}]),t}(c.a.Component),le={checkLogin:function(){return function(){var e=Object(q.a)(H.a.mark(function e(t){return H.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t($()),e.abrupt("return",new Promise(function(e){var n=setTimeout(function(){return Promise.resolve()},3e3);Promise.race([n,function(){return new Promise(function(e){for(;!window.gapi;);e()})}]),window.gapi&&gapi.load("auth2",function(){gapi.auth2.init({client_id:"634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com"}).then(function(){if(gapi.auth2.getAuthInstance().isSignedIn.get()){var e=gapi.auth2.getAuthInstance().currentUser.get().getBasicProfile();t(ee(e.getName(),e.getName(),e.getEmail()))}}),e()})}));case 2:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}()},openConnection:function(){return function(e){e({type:L});var t=new WebSocket("ws://".concat("127.0.0.1:3001","/ws"));t.onopen=function(){t.send(JSON.stringify({type:0,message:"Hi From the Client! The websocket just opened"})),t.send(JSON.stringify({type:1,userId:"AAAAAA",x:6,y:9,r:4,g:2,b:0})),e(function(e){return{type:_,payload:{socket:e}}}(t))},t.onclose=function(t){e(X())},t.onerror=function(t){e(function(e){return{type:A,payload:{error:e}}}(t.type))},t.onmessage=function(e){var t=e.data;console.log("Recieved a message from the server, message: "+t)}}}},se=Object(v.b)(null,le)(ue),fe=n(20),pe=n(8);function me(e,t){return function(){var n=arguments.length>0&&void 0!==arguments[0]?arguments[0]:e,a=arguments.length>1?arguments[1]:void 0,r=t[a&&a.type];return r?r(n,a):n}}var de,ge,ve,be,Oe,he,Ee,ye=me(null,(a={},Object(pe.a)(a,A,function(){return null}),Object(pe.a)(a,"canvas/REGISTER_CONTEXT",function(e,t){return t.payload.context}),a)),je=me(null,(r={},Object(pe.a)(r,"canvas/CLEAR_POSITION",function(){return null}),Object(pe.a)(r,"canvas/UPDATE_POSITION",function(e,t){return t.payload}),r)),Ce=me(41,Object(pe.a)({},"canvas/SET_ZOOM",function(e,t){return t.payload.zoom})),ke=Object(fe.c)({canvasContext:ye,curPosition:je,zoomFactor:Ce}),we=me(!1,(de={},Object(pe.a)(de,V,function(){return!1}),Object(pe.a)(de,K,function(){return!1}),Object(pe.a)(de,Q,function(){return!1}),Object(pe.a)(de,W,function(){return!0}),de)),Se=me(null,(ge={},Object(pe.a)(ge,V,function(){return null}),Object(pe.a)(ge,K,function(){return null}),Object(pe.a)(ge,Q,function(){return null}),Object(pe.a)(ge,W,function(e,t){return t.payload.name}),ge)),Ne=me(null,(ve={},Object(pe.a)(ve,V,function(){return null}),Object(pe.a)(ve,K,function(){return null}),Object(pe.a)(ve,Q,function(){return null}),Object(pe.a)(ve,W,function(e,t){return t.payload.email}),ve)),Pe=me(null,(be={},Object(pe.a)(be,V,function(){return null}),Object(pe.a)(be,K,function(){return null}),Object(pe.a)(be,Q,function(){return null}),Object(pe.a)(be,W,function(e,t){return t.payload.id}),be)),Re=Object(fe.c)({isLoggedIn:we,name:Se,email:Ne,id:Pe}),Te=me(null,(Oe={},Object(pe.a)(Oe,L,function(){return null}),Object(pe.a)(Oe,A,function(){return null}),Object(pe.a)(Oe,M,function(){return null}),Object(pe.a)(Oe,_,function(e,t){return t.payload.socket}),Oe)),xe=me(!1,(he={},Object(pe.a)(he,L,function(){return!1}),Object(pe.a)(he,A,function(){return!1}),Object(pe.a)(he,M,function(){return!1}),Object(pe.a)(he,_,function(){return!0}),he)),Ie=me(!1,(Ee={},Object(pe.a)(Ee,L,function(){return!1}),Object(pe.a)(Ee,A,function(){return!0}),Object(pe.a)(Ee,M,function(){return!1}),Object(pe.a)(Ee,_,function(){return!1}),Ee)),Le=Object(fe.c)({socket:Te,isConnected:xe,receivedError:Ie}),_e=n(202);n.d(t,"store",function(){return Xe});var Ae=Object(fe.c)({canvas:ke,login:Re,websocket:Le}),Me=window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__||fe.d,Xe=Object(fe.e)(Ae,Me(Object(fe.a)(_e.a)));Object(i.render)(c.a.createElement(v.a,{store:Xe},c.a.createElement(se,null)),document.getElementById("root"))}},[[205,1,2]]]);
//# sourceMappingURL=main.f5ba173d.chunk.js.map