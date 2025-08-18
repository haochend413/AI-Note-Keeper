"""
LangGraph based multi-agent chatbot system with persistent input.
"""

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from langchain_ollama import ChatOllama
from langchain_mcp_adapters.tools import load_mcp_tools
from langgraph.prebuilt import create_react_agent
from langgraph_supervisor import create_supervisor
from pathlib import Path
from langchain.schema import HumanMessage, SystemMessage
from langchain.tools import tool
from langchain.agents import Tool

import asyncio

server_params = StdioServerParameters(
    command="python",
    args=[str(Path(__file__).parent.parent / "server" / "server.py")],
)


async def main():
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()

            # Load MCP tools (data fetching)
            tools = await load_mcp_tools(session)

            # Initialize LLM
            llm = ChatOllama(model="llama3.1:latest")

            # Tool: summarize_state
            @tool
            def summarize_state(messages: list) -> str:
                """Summarize fetched notes and topics into goals, tasks, and identity."""
                system_msg = SystemMessage(
                    content=(
                        "You are a summary agent. Summarize the fetched notes and topics, "
                        "and conclude with what are my main goals, what I have to do, and who I am."
                    )
                )
                response = llm.invoke([system_msg] + messages)
                return response.content

            summary_tool = Tool(
                name="summarize_state",
                func=summarize_state,
                description="Summarizes fetched notes and topics into main goals, tasks, and who the user is.",
            )

            # Data fetching agent
            read_agent = create_react_agent(
                model=llm,
                tools=tools,
                prompt="You are a data fetching agent. Your job is to fetch all the notes or topics of a given range. YOU MUST CALL THE TOOLS PROVIDED TO YOU.",
                name="read_agent",
                # verbose=True,
                debug=True,
            )

            # Summary agent
            summary_agent = create_react_agent(
                model=llm,
                tools=[summary_tool],
                prompt="You are a summary agent. Use the summarize_state tool to summarize notes.",
                name="summary_agent",
            )

            # Supervisor agent
            supervisor = create_supervisor(
                agents=[read_agent, summary_agent],
                model=llm,
                prompt=(
                    "You are a supervisor. Coordinate between assistants by delegating tasks ONLY "
                    "through their tools. First ask the data fetching assistant to gather notes, "
                    "then pass results to the summary assistant to summarize. Never answer directly yourself."
                ),
            ).compile()

            print("✅ Chatbot running. Type your message (Ctrl+C to exit).")

            # Persistent chat loop
            while True:
                user_input = input("\nYou: ").strip()
                if not user_input:
                    continue

                # Stream AI responses
                for chunk in supervisor.stream(
                    {"messages": [{"role": "user", "content": user_input}]}
                ):
                    # # Only print AI messages
                    # if chunk.get("type") in ("AIMessage", "tool_message"):
                    #     content = chunk.get("content")
                    #     if content:
                    #         print(f"\nAI: {content}\n", end="")
                    print(chunk)


# Run the chat loop
asyncio.run(main())
