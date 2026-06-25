---
name: Rspack webpackContext regex
description: import.meta.webpackContext keys in Rspack include "./" prefix — regex must account for it or no modules load
---

## Rule
When using `import.meta.webpackContext` in Rsbuild/Rspack, all keys returned by `.keys()` are prefixed with `./` (e.g., `./api.ts`, not `api.ts`). The regex passed to webpackContext is tested against these prefixed keys.

**Why:** The original router used `/^[^.]+\.ts$/` which requires the string to start with a non-dot character. Since all keys start with `./` (a dot), zero keys matched and `menuList` was empty — all Vue Router routes silently disappeared.

**How to apply:** Always start the regex with `\./` when filtering webpackContext module keys:
```js
// WRONG (fails in Rspack — all routes disappear)
regExp: /^[^.]+\.ts$/

// CORRECT (works in Rspack)
regExp: /^\.\/[^.]+\.ts$/
```

The `[^.]+` part still correctly excludes test files: `./api.test.ts` fails because after matching `./api`, the next char is `.` which stops `[^.]+`, then `\.ts` matches `.te` but then `$` fails since `.ts` remains.

**Also:** After fixing the regex, also guard against both ES module format (`mod.default`) and CommonJS format (direct `mod` object) when extracting the route:
```js
const route = (mod?.default?.path ? mod.default : undefined)
           ?? ((mod as RouteRecordRaw)?.path ? (mod as RouteRecordRaw) : undefined)
```
