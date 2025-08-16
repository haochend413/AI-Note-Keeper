# Create server parameters for stdio connection
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from langchain_ollama import ChatOllama
from langchain_mcp_adapters.tools import load_mcp_tools
from langgraph.prebuilt import create_react_agent
from pathlib import Path
import asyncio

server_params = StdioServerParameters(
    command="python",
    # Make sure to update to the full absolute path to your math_server.py file
    args=[str(Path(__file__).parent.parent / "server" / "server.py")],
)


async def main():
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize the connection
            await session.initialize()

            # Get tools
            tools = await load_mcp_tools(session)
            # Initialize Ollama LLM
            llm = ChatOllama(model="mistral")

            # Bind tools to LLM if required
            # llm_with_tools = llm.bind_tools(tools)  # optional depending on adapter

            # Create React agent
            agent = create_react_agent(llm, tools)
            agent_response = await agent.ainvoke(
                {
                    "messages": [
                        {
                            "role": "user",
                            "content": "Read all my notes and give me a summary",
                        }
                    ]
                }
            )
            print(agent_response)


asyncio.run(main())
