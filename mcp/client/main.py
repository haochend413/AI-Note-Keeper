from client import MCPClient
from langchain.agents import initialize_agent
import asyncio
from langchain_community.llms import Ollama

from pathlib import Path
from langchain.agents import initialize_agent, AgentType


async def main():
    client = MCPClient()
    await client.connect_to_server(
        str(Path(__file__).parent.parent / "server" / "server.py")
    )

    try:
        tools = await client.get_mcp_tools()

        agent = initialize_agent(
            tools,
            llm=Ollama(model="mistral"),
            agent=AgentType.CHAT_ZERO_SHOT_REACT_DESCRIPTION,
            verbose=True,
            handle_parsing_errors=True,
        )

        # Async execution
        result = await agent.arun("List all notes in the database that are not deleted")
        print(result)
    finally:
        # Ensure we close the connection properly
        await client.close()


if __name__ == "__main__":
    asyncio.run(main())
