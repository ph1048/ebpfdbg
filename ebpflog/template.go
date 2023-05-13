package ebpflog

import "fmt"

var header string = `
<html>
<head>
<style>
body {
	font-family: Consolas;
	background: #eaf7fa;
}
.src {
	color:#833fcc;
}
.bpf {
	color:#139f7c;
}
.err {
	color:#ff0026;
	font-weight:bold;
}
.bt {
	color:#a0a362;
}
.result {
	color:#009bb9;
}
.bttbl {
	background:#3bf1c4;
}
.regdump {
	background:#faf7c2;
	display:inline-block;
}
.backtrace {
	background:#c2f0fa;
	display:inline-block;
}
.regdumpcontent {
	display:none;
}
.registers {
	display:inline-block;
}
.registerst td, .registerst th {
	border:1px solid black;
	margin:0;
	background:#fbfae8;
}
.registerst th {
	font-weight:bold;
}
.stack {
	display:inline-block;
}
.btinfo {
	border:1px solid green;
}
.pclnk, .pclnk:visited {
	color: #3371db;
}
.icon {
	text-decoration: none;
	font-family: "Noto Color Emoji", "Apple Color Emoji", "Segoe UI Emoji",	Times, Symbola, Aegyptus, Code2000, Code2001, Code2002, Musica, serif, LastResort;
}
.button {
	cursor: pointer;
	text-decoration: underline dotted;
}
.annotation {
	padding-left:10px;
}
.starting {
	height: 50px;
	border-bottom: 1px dashed gray;
}
.ending {
	height: 50px;
	border-top: 3px dashed #04F013;
	margin-top: 20px;
}
.bbtrans {
	font-size:40px;
}
</style>
<script>
function toggleshow(e) {
	x = e.parentElement.getElementsByClassName("regdumpcontent")[0]
	if (x.style.display == "none" || x.style.display == "") {
		x.style.display = "block";
	  } else {
		x.style.display = "none";
	  }
}
var lasthlidx = 0
function onregclk(e) {
	if (lasthlidx != 0) {
		document.styleSheets[0].deleteRule(lasthlidx)
		lasthlidx = 0
	}
	if (e == null) {
		return
	}
	
	tgt = e.target
	cl = tgt.className
	lasthlidx = document.styleSheets[0].cssRules.length
	document.styleSheets[0].addRule("."+cl, "background:#3b161f;color:#a4f32c;")
}

window.onload = function() {
	var anchors = document.getElementsByTagName('span');
	for(var i = 0; i < anchors.length; i++) {
		var anchor = anchors[i];
		if(anchor.className.startsWith("reg_")) {
			anchor.onclick = function(e) {
				onregclk(e);
			}
		}
	}
}
</script>
</head>
<body>
`
var footer string = `
</body>
</html>
`

func getIcon(icon string) string {
	return fmt.Sprintf("<span class=\"icon\">%s</span>", icon)
}
