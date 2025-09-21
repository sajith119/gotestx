# 🚀 GoTestX v1.1.0 Release Notes

GoTestX extends `go test` with coverage support, quiet mode, and clean output filtering.  
This release improves **quiet mode behavior**, renames the clean flag, and strengthens test coverage.

---

## ✨ What’s New
- **Quiet mode summary output**  
  - Always prints a concise result instead of nothing:
    - ✅ success if all tests passed  
    - coverage % if `-c` is enabled  
    - ❌ failure with hint to rerun without `-q`  

- **New flag `-V, --clean-view`**  
  - Replaces the old `-C` flag.  
  - Suppresses `? ... [no test files]` lines for a cleaner view.  
  - Works alongside `-q` and `-c`.  

- **Unified error messages in quiet mode**  
  - All failures now print:  
    ```
    ❌ Tests failed (use without -q to see details)
    ```

---

## 🛠 Improvements
- Stronger **test coverage**:
  - Quiet mode failure path covered.
  - Default `./...` package expansion verified.
  - Direct ellipsis (`./...`) handling tested.

- **Documentation updates**:
  - Added **Sample Output** section to README.
  - Updated usage to reflect `-V` flag.

---

## 📦 Install / Upgrade
```bash
go install github.com/entiqon/gotestx@latest
```

Check version:
```bash
gotestx -v
```

---

🔖 Tag: `v1.1.0`  
📄 License: MIT  
