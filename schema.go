package main

const RuleSetSchema = `{
    "name": "Blueprint Rules",
    "description": "Default rule configuration for Genesys Cloud Blueprints",
    "ruleGroups": {
        "STRUCT": {
            "description": "Validation for required file/folder existence",
            "rules": {
                "1": {
                    "description": "All Genesys Cloud blueprints must include a README.MD file.  This file should contain a brief introduction of the blueprint.",
                    "path": "./README.md",
                    "conditions": [{ "pathExists": true }],
                    "level": "error"
                },
                "2": {
                    "description": "Every Genesys Cloud blueprint should have a blueprint directory at the root of the project.  This directory should hold all assets associated with the blueprint.",
                    "path": "./blueprint",
                    "conditions": [{ "pathExists": true }],
                    "level": "error"
                },
                "3": {
                    "description": "Every Genesys Cloud blueprint should have a blueprint/images directory that will contain all of the image assets for a project.",
                    "path": "./blueprint/images",
                    "conditions": [{ "pathExists": true }],
                    "level": "error"
                },
                "4": {
                    "description": "Every Genesys Cloud blueprint should have a blueprint/images directory that will contain all of the image assets for a project.",
                    "path": "./blueprint/images",
                    "conditions": [{ "pathExists": true }],
                    "level": "error"
                },
                "5": {
                    "description": "Every Genesys Cloud blueprint should have a blueprint/index.md that contains a complete writeup in Markdown of the blueprint.",
                    "path": "./blueprint/index.md",
                    "conditions": [{ "pathExists": true }],
                    "level": "error"
                }
            }
        },
        "CONTENT": {
            "description": "Content related validation",
            "rules": {
                "1": {
                    "description": "Overview image should be referred to in README.MD",
                    "path": "./README.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "!\\[.*\\]\\(images/overview\\.png +['|\"].*['|\"]\\)"
                        }]
                    }],
                    "level": "error"
                },
                "2": {
                    "description": "The front matter must be defined in the file or the blueprint will not appear in the Developer Center",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "^---\\r*\\n*.*\\r*\n*---"
                        }]
                    }],
                    "level": "error"
                },
                "3": {
                    "description": "  The front matter must be defined in the file or the blueprint will not appear in the Developer Center.",
                    "path": "./README.md",
                    "conditions": [{
                        "markdownMeta": {
                            "title": ".*",
                            "author": ".*",
                            "indextype": "blueprint",
                            "icon": "blueprint",
                            "image": ".*",
                            "category": ".*",
                            "summary": ".*"
                        }
                    }],
                    "level": "error"
                },
                "4": {
                    "description": "The index.md must have a ## Scenario section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Scenario *"
                        }]
                    }],
                    "level": "error"
                },
                "5": {
                    "description": "The index.md must have a ## Solution section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Solution *"
                        }]
                    }],
                    "level": "error"
                },
                "6": {
                    "description": "The index.md must have a ## Content section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Content *"
                        }]
                    }],
                    "level": "error"
                },
                "7": {
                    "description": "The index.md must have a ## Prerequisites section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Prerequisites *"
                        }]
                    }],
                    "level": "error"
                },
                "8": {
                    "description": "The index.md must have a ### Specialized knowledge section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "### *Specialized knowledge *"
                        }]
                    }],
                    "level": "error"
                },
                "9": {
                    "description": "The index.md must have a ## Implementation steps section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Implementation steps *"
                        }]
                    }],
                    "level": "error"
                },
                "10": {
                    "description": "The index.md must have a ### Download the repository containing the project files section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "### *Download the repository containing the project files *"
                        }]
                    }],
                    "level": "error"
                },
                "11": {
                    "description": "The index.md must have a ## Additional resources section describing the problem the blueprint is trying to solve.",
                    "path": "./blueprint/index.md",
                    "conditions": [{
                        "contains": [{
                            "type": "regex",
                            "value": "## *Additional resources *"
                        }]
                    }],
                    "level": "error"
                }
            }
        },
        "LINK": {
            "description": "Validates the links in Markdown files",
            "rules": {
                "1": {
                    "description": "Image links in the README.md or index.md file should point to a valid image file.",
                    "files": ["./readme.md", "./blueprint/index.md"],
                    "conditions": [{
                        "checkReferenceExist": ["!\\[.*\\]\\((.*\\) *\".*\").*"]
                    }],
                    "level": "error"
                },
                "2": {
                    "description": "Image links in the README.md or index.md file is missing alternative text.",
                    "files": ["./readme.md", "./blueprint/index.md"],
                    "conditions": [{
                        "notContains": ["!\\[.*\\]\\(.*[^ \"]+[^\"]*\\)"]
                    }],
                    "level": "error"
                },
                "3": {
                    "description": "Hyperlinks in the README.md or index.md file is missing alternative text.",
                    "files": ["./readme.md", "./blueprint/index.md"],
                    "conditions": [{
                        "notContains": ["\\[.*\\]\\(.*[^ \"]+[^\"]*\\)"]
                    }],
                    "level": "error"
                }
            }
        }
    }
}

`
