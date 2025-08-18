from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from langchain_ollama import ChatOllama
from langchain_mcp_adapters.tools import load_mcp_tools
from langgraph.prebuilt import create_react_agent
from pathlib import Path
import asyncio
import json
from rich.console import Console
from rich.panel import Panel
from rich.markdown import Markdown
from rich.table import Table
import re
from langchain.schema import HumanMessage, AIMessage, SystemMessage
from langchain_core.messages import ToolMessage  # Add this import

# Initialize Rich console
console = Console()

# Create server parameters for stdio connection
server_params = StdioServerParameters(
    command="python",
    args=[str(Path(__file__).parent.parent / "server" / "server.py")],
)


def extract_tool_calls(content):
    """Extract tool calls from response content."""
    if not isinstance(content, str):
        return []

    tool_pattern = r'{"name": "([^"]+)", "parameters": ({[^}]+})}'
    matches = re.findall(tool_pattern, content)

    tool_calls = []
    for tool_name, params_str in matches:
        try:
            params = json.loads(params_str)
            tool_calls.append({"name": tool_name, "parameters": params})
        except json.JSONDecodeError:
            tool_calls.append({"name": tool_name, "parameters": params_str})

    return tool_calls


async def main():
    # Show startup banner
    console.print(
        Panel.fit(
            "[bold blue]Mantis Notes Assistant[/bold blue]\n"
            "[italic]Powered by LangGraph + Ollama[/italic]",
            border_style="blue",
        )
    )

    console.print("\n[yellow]Connecting to server...[/yellow]")

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize the connection
            await session.initialize()
            console.print("[green]✓[/green] Server connection established")

            # Get tools
            tools = await load_mcp_tools(session)

            # Display available tools
            tool_table = Table(title="Available Tools")
            tool_table.add_column("Tool Name", style="cyan")
            tool_table.add_column("Description", style="green")

            for tool in tools:
                tool_table.add_row(tool.name, tool.description)

            console.print(tool_table)
            console.print(f"[green]✓[/green] Loaded {len(tools)} tools from server")

            # Initialize Ollama LLM
            llm = ChatOllama(model="llama3.1:latest")
            console.print("[green]✓[/green] Connected to Ollama LLM")

            # Create React agent
            agent = create_react_agent(llm, tools)
            console.print("[green]✓[/green] Agent initialized and ready")

            # Print usage instructions
            console.print("\n[bold]Enter your questions about your notes below.[/bold]")
            console.print("[dim]Type 'exit' to quit the program.[/dim]\n")

            # Create the chat loop
            while True:
                # Get user input
                user_input = console.input("\n[bold cyan]You: [/bold cyan]").strip()
                if not user_input:
                    continue
                if user_input.lower() == "exit":
                    console.print("[yellow]Exiting...[/yellow]")
                    break

                # Process the input
                console.print("\n[cyan]Thinking...[/cyan]")

                try:
                    resp = await agent.ainvoke(
                        {"messages": [{"role": "user", "content": user_input}]}
                    )

                    # Extract and display tool calls
                    messages = resp.get("messages", [])

                    # Tracking for tool calls and final answer
                    tool_calls_found = False
                    final_answer = None

                    # Process each message
                    for i, message in enumerate(messages):
                        # Handle different message types
                        if isinstance(message, dict):
                            role = message.get("role", "")
                            content = message.get("content", "")
                        elif isinstance(
                            message, (HumanMessage, AIMessage, SystemMessage)
                        ):
                            # LangChain message objects
                            role = message.type
                            content = message.content
                        elif isinstance(message, ToolMessage):
                            # Handle ToolMessage
                            role = "tool"
                            content = message.content

                            # Display tool result
                            console.print(
                                Panel(
                                    f"[bold yellow]Tool ID:[/bold yellow] {message.tool_call_id}\n"
                                    f"[bold yellow]Result:[/bold yellow] {content}",
                                    title="[bold magenta]Tool Result[/bold magenta]",
                                    border_style="magenta",
                                )
                            )
                            continue  # Skip regular processing for tool messages
                        else:
                            # Unknown message type
                            console.print(
                                f"[yellow]Unknown message type: {type(message)}[/yellow]"
                            )
                            continue

                        # Skip non-assistant messages for tool extraction
                        if role != "assistant" and role != "ai":
                            continue

                        # Try to extract tool calls
                        tool_calls = extract_tool_calls(content)

                        if tool_calls:
                            tool_calls_found = True
                            tool_call_table = Table(title=f"Tool Call #{i+1}")
                            tool_call_table.add_column("Tool", style="magenta")
                            tool_call_table.add_column("Parameters", style="yellow")

                            for call in tool_calls:
                                tool_call_table.add_row(
                                    call["name"],
                                    json.dumps(call["parameters"], indent=2),
                                )

                            console.print(tool_call_table)
                        elif role == "assistant" or role == "ai":
                            # This could be the final answer
                            final_answer = content

                    # Show final answer (the last non-tool assistant message)
                    if final_answer:
                        console.print(
                            Panel(
                                Markdown(final_answer),
                                title="[bold green]Assistant[/bold green]",
                                border_style="green",
                            )
                        )
                    else:
                        console.print("[red]No clear answer was generated.[/red]")

                    # If we didn't find any tool calls but have messages, show the raw response
                    if not tool_calls_found and not final_answer and messages:
                        console.print(
                            "[yellow]No structured response found. Raw output:[/yellow]"
                        )
                        console.print(resp)

                except Exception as e:
                    console.print(f"[bold red]Error:[/bold red] {str(e)}")
                    import traceback

                    console.print(traceback.format_exc())


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        console.print("\n[yellow]Program terminated by user[/yellow]")
    except Exception as e:
        console.print(f"\n[bold red]Fatal error:[/bold red] {str(e)}")
