import type { Argv } from "yargs"
import { runTask } from "../../agents/planLoop"
import { UI } from "../ui"
import { cmd } from "./cmd"

export const PlanCommand = cmd({
  command: "plan <goal>",
  describe: "run the two-agent plan-step-review loop",
  builder: (yargs: Argv) => {
    return yargs
      .positional("goal", {
        describe: "the task goal to accomplish",
        type: "string",
        demandOption: true,
      })
      .option("dry-run", {
        describe: "show the plan without executing steps",
        type: "boolean",
        default: false,
      })
  },
  handler: async (argv) => {
    const { goal, dryRun } = argv

    UI.println(`${UI.Style.TEXT_INFO_BOLD}Planning task:${UI.Style.TEXT_NORMAL} ${goal}`)

    if (dryRun) {
      UI.println(
        `${UI.Style.TEXT_WARNING_BOLD}Dry run mode${UI.Style.TEXT_NORMAL} - plan will be displayed but not executed`,
      )
      // In dry run mode, we could show the plan without executing
      UI.println(`${UI.Style.TEXT_WARNING}Dry run not yet implemented${UI.Style.TEXT_NORMAL}`)
      return
    }

    try {
      await runTask(goal as string)
      UI.println(`${UI.Style.TEXT_SUCCESS_BOLD}Task completed successfully!${UI.Style.TEXT_NORMAL}`)
    } catch (error) {
      UI.error(`Task failed: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})
