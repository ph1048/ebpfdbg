# ebpfdbg (eBPF verifier log viewer)
This project is aimed to help debug large eBPF verifier error logs.
This simple utility allows to post-process eBPF verifier log to a human-readable HTML page.


一次滂沱大雨的来临，你需要一把保护伞。eBPF排错程序为你的软件排忧解难。

---

## Quick usage guide
### 1. Get a file with eBPF verifier log
If you see an error message from eBPF verifier due to loading of your program, you need to extract full eBPF verifier logs.
Make sure your log is not truncated to default 65535 bytes. Save it to a file.
### 2. Run eBPF verifier against it
You need to have Go compiler 1.19 or higher on your system.
Run the following:

```
go run github.com/ph1048/ebpfdbg/cmd/ebpfdbg@v1.0.2 serve --input verifier.log
```

### 3. Open URL in web browser
Depending on the verifier log size, this page might be heavy.

---

## Screenshots

![screenshot](pic/scr1.png "Screenshot")
![screenshot](pic/scr2.png "Screenshot")

## Contribution
Ways to contribute:
1. Create issues for problems and suggestions
2. Post your full eBPF verifier logs (if possible)

