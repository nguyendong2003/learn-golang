# Pixel-Perfect React Native Code Generation Plan from Figma Designs

This plan outlines the creation of a set of agent skills in the `.agents/skills` folder to automate the generation of pixel-perfect React Native screens and components (using Expo Router and `StyleSheet`) directly from the JSON files exported by `test_figma_mcp.py` via `figma-mcp-go`.

## Goal Description

The project requires converting Figma designs into highly accurate (pixel-perfect) React Native code. To achieve this, the user has created `test_figma_mcp.py` to extract selections, styles, and screenshots from Figma. The selection JSON is sorted vertically and horizontally (handling sub-pixel layout issues and centering alignment).

We will create a set of agent skill files (in **English**) in `.agents/skills/` to define the screen/component building process, naming conventions, asset/layout parsing, style extraction, and Git workflow rules.

---

## User Review Required

Please review the revised list of skill files, focusing on the component storage path:

> [!IMPORTANT]
> All generated components, whether common, layout-specific, or sub-components, must be stored under **`src/components/`** directly. We will **not** create or use a `src/components/ui/` folder.

---

## Proposed Changes

### Agent Skills Folder

#### [NEW] [fe-gen-screen.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/fe-gen-screen.md)
Coordinates screen-level generation:
1. Run/inspect the Figma export (JSON files, screenshots) for the screen.
2. Scan for and export new assets (using `asset_manager.md`).
3. Extract and map shared styles (using `style_mapper.md`).
4. **Decomposition**: Break down the screen layout into smaller, logical sub-components.
5. **Component Reuse Check**: Scan `src/components/` for similar components. If one exists, reuse it or extend its props/style.
6. Generate the layout with React Native `StyleSheet` (using `parse_layout.md`), referencing the shared styles and sub-components.
7. Apply naming conventions and verify code.

#### [NEW] [fe-gen-component.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/fe-gen-component.md)
Coordinates standalone component generation (when selecting a component node in Figma):
1. Run/inspect the Figma export (JSON, screenshot) for the component.
2. Scan for and export new assets (using `asset_manager.md`).
3. **Reuse Check**: Compare with existing components. If a similar one is found, modify/extend the existing component (e.g., adding props, style customization, extra options) instead of writing a new one.
4. If a new component is necessary, decompose it if it has nested complex elements.
5. Generate the component code directly under `src/components/` with TypeScript props using `parse_layout.md`.
6. Export the component for general use.

#### [NEW] [naming_conventions.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/naming_conventions.md)
Establishes clear naming rules across the project:
- **Folders**: lowercase kebab-case (e.g. `src/components/common-inputs`).
- **Files**: PascalCase for React Native components/screens (e.g. `ExploreHeader.tsx`), kebab-case or camelCase for utility/helper files.
- **Variables/Constants**: camelCase for normal variables, UPPER_CASE for static configuration constants.
- **Styles**: camelCase (e.g. `containerStyle`, `headingText`).
- **Functions**: camelCase starting with a verb (e.g. `handleSubmit`, `renderItem`).

#### [NEW] [git_guidelines.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/git_guidelines.md)
Guidelines for standardizing source control contributions:
- **Branch Naming**: `feature/screen-name`, `bugfix/issue-name`, `refactor/component-name`.
- **Commit Message Format**: Following Conventional Commits (`feat(screen-name): add pixel-perfect layout`, `fix(icons): resolve layout offset`, `style(theme): update primary colors`).

#### [NEW] [parse_layout.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/parse_layout.md)
Defines how to translate geometric bounds and hierarchy from the Figma JSON tree into a clean React Native Flexbox structure. It handles:
- Flex direction (Row vs. Column detection based on bounds alignment).
- Padding, margins, gap calculations.
- Text style mapping (font family, font weight mapping from word labels to React Native format, line-height conversion).
- Absolute positioning detection for overlapping components (like badges, floating buttons, background shapes).

#### [NEW] [style_mapper.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/style_mapper.md)
Instructs how to parse the JSON output of the `get_styles` tool (common Figma styles) and translate it into a structured theme file, e.g., in `src/constants/theme.ts`. It covers:
- Extracting color palettes (primary, secondary, neutrals).
- Mapping typography scale (body, headers, titles).
- Replacing hardcoded styling in generated code with theme variables (e.g., `theme.colors.primary` instead of `"#3629b7"`).

#### [NEW] [asset_manager.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/asset_manager.md)
Automates asset generation:
- Scanning JSON nodes to detect images (`RECTANGLE` with image fills, or `VECTOR` paths) and icons (node names containing `icon`, `logo`, or shape-based indicators).
- Calling the `figma-mcp-go` tool `save_screenshots` to export them directly to `assets/icons/` (in `SVG` format) or `assets/images/` (in `PNG` format).
- Creating/updating registry files (`assets/icons/index.ts`, `assets/images/index.ts`) for easy clean imports.
- Recommending `expo-image` as the renderer for performance and ease-of-use (SVG & PNG support).

#### [NEW] [component_extractor.md](file:///home/gmo/Documents/Learning/LearnJS/react-native-figma-test/.agents/skills/component_extractor.md)
Controls modular component extraction:
- Detecting duplicate nodes or repeated visual structures across screen JSON files or inside a single screen (e.g., custom input fields, primary buttons, headers).
- Generating reusable components directly in `src/components/` with TypeScript props (e.g. `label`, `value`, `onChangeText`, `onPress`).
- Ensuring the screen files import these common components instead of writing inline code.

---

## Verification Plan

### Automated Verification
- Run a TypeScript build check `npx tsc --noEmit` and the lint command `npm run lint` if assets and files are generated, verifying that the structure and types are completely correct.
- Verify the skill files exist in `.agents/skills` and contain clear markdown headings and structured English instructions.

### Manual Verification
- Verify that all eight skills are read and properly interpreted.
