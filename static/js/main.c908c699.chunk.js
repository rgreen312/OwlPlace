(window.webpackJsonpclient=window.webpackJsonpclient||[]).push([[0],{115:function(e,n,t){e.exports=t(187)},120:function(e,n,t){},121:function(e,n,t){},186:function(e,n,t){},187:function(e,n,t){"use strict";t.r(n);var a=t(0),o=t.n(a),c=t(5),l=(t(120),t(121),t(49)),r=t(30),i=Object(r.e)(function(e){var n=e.children,t=e.location.pathname;return Object(a.useEffect)(function(){window.scrollTo(0,0)},[t]),n||null}),u=t(43),s=function(e){e.onClick;return o.a.createElement("div",null,o.a.createElement("div",null,o.a.createElement("canvas",null)))},m={onClick:function(){return function(e,n){console.log("dispatching action"),e({type:"canvas/FETCH_IMAGE_START",payload:{data:"testing"}})}}},g=Object(u.b)(null,m)(s),f=function(){return o.a.createElement(g,null)},d=function(){return o.a.createElement("div",{className:"about-page"},"Update the canvas one pixel at a time...")},p=t(61),E=t(85),b=t(190),O=t(10),v=t(189),h=function(e){var n=e.onLogin,t=e.isLoggedIn,a=e.name,c=t?o.a.createElement(o.a.Fragment,null,"Hi, ",a):o.a.createElement(p.a,{className:"login-button",onClick:n},"Login"),r=o.a.createElement(E.a,null,o.a.createElement(E.a.Item,null,o.a.createElement(l.b,{to:"/about"},"About")),t&&o.a.createElement(E.a.Item,null,o.a.createElement("button",{onClick:function(){}},"Sign out"))),i=o.a.createElement(b.a,{key:"more",overlay:r},o.a.createElement(p.a,{style:{border:"none",padding:0}},o.a.createElement(O.a,{type:"ellipsis",style:{fontSize:20,verticalAlign:"top"}})));return o.a.createElement(v.a,{title:"OwlPlace",extra:[c,i]})},j=t(84),y=t.n(j),w=t(112),C="login/LOGIN_START",I="login/LOGIN_SUCCESS",S="login/LOGIN_ERROR",k=function(e){return e.login.isLoggedIn},_=function(e){return e.login.name||void 0},L={onLogin:function(){return function(){var e=Object(w.a)(y.a.mark(function e(n){var t,a,o;return y.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return n({type:C}),t=new Promise(function(e){gapi.load("auth2",function(){gapi.auth2.init({client_id:"634069824484-ch6gklc2fevg9852aohe6sv2ctq7icbk.apps.googleusercontent.com"}),e()})}),e.next=4,t;case 4:return e.next=6,gapi.auth2.getAuthInstance().signIn();case 6:a=gapi.auth2.getAuthInstance().currentUser.get(),o=a.getBasicProfile(),n((c=o.getName(),l=o.getId(),r=o.getEmail(),{type:I,payload:{name:c,id:l,email:r}}));case 9:case"end":return e.stop()}var c,l,r},e)}));return function(n){return e.apply(this,arguments)}}()}},N=Object(u.b)(function(e){return{isLoggedIn:k(e),name:_(e)}},L)(h),T=(t(186),function(){return o.a.createElement("div",null,o.a.createElement("footer",null,"Made with \u2665 by COMP 413 @ Rice"))}),A=function(){return o.a.createElement(l.a,null,o.a.createElement(i,null,o.a.createElement(N,null),o.a.createElement(r.a,{exact:!0,path:"/",component:f}),o.a.createElement(r.a,{path:"/about",component:d}),o.a.createElement(T,null)))},R=new WebSocket("ws://127.0.0.1:3010/ws");console.log("Attempting Connection..."),R.onopen=function(){console.log("Successfully Connected"),R.send(JSON.stringify({type:0,message:"Hi From the Client! The websocket just opened"}))},R.onclose=function(e){console.log("Socket Closed Connection: ",e),R.send(JSON.stringify({type:9,message:"Client Closed!"}))},R.onerror=function(e){console.log("Socket Error: ",e)},R.onmessage=function(e){var n=e.data;console.log("Recieved a message from the server, message: "+n)};var x,P,G,J,M=function(){return o.a.createElement(A,null)},U=t(22),F=t(15);function H(e,n){return function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:e,a=arguments.length>1?arguments[1]:void 0,o=n[a&&a.type];return o?o(t,a):t}}var B=H(!1,(x={},Object(F.a)(x,C,function(){return!1}),Object(F.a)(x,S,function(){return!1}),Object(F.a)(x,I,function(){return!0}),x)),D=H(null,(P={},Object(F.a)(P,C,function(){return null}),Object(F.a)(P,S,function(){return null}),Object(F.a)(P,I,function(e,n){return n.payload.name}),P)),X=H(null,(G={},Object(F.a)(G,C,function(){return null}),Object(F.a)(G,S,function(){return null}),Object(F.a)(G,I,function(e,n){return n.payload.email}),G)),q=H(null,(J={},Object(F.a)(J,C,function(){return null}),Object(F.a)(J,S,function(){return null}),Object(F.a)(J,I,function(e,n){return n.payload.id}),J)),z=Object(U.c)({isLoggedIn:B,name:D,email:X,id:q}),V=t(113),W=Object(U.c)({login:z}),K=window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__||U.d,Q=Object(U.e)(W,K(Object(U.a)(V.a)));Object(c.render)(o.a.createElement(u.a,{store:Q},o.a.createElement(M,null)),document.getElementById("root"))}},[[115,1,2]]]);
//# sourceMappingURL=main.c908c699.chunk.js.map