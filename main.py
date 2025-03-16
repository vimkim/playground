import curses
import random
import time

TICK_RATE = 0.2  # 200 milliseconds per tick
NEW_BLOCK_PROBABILITY = 0.1  # probability per column per tick


class FallingBlock:
    def __init__(self, x, y):
        self.x = x
        self.y = y


class Game:
    def __init__(self, width, height):
        self.player_x = width // 2
        self.player_y = height - 1
        self.blocks = []
        self.score = 0
        self.width = width
        self.height = height

    def update(self):
        # Spawn new blocks along the top row for each column.
        for x in range(self.width):
            if random.random() < NEW_BLOCK_PROBABILITY:
                self.blocks.append(FallingBlock(x, 0))
        # Move blocks down and remove those that have moved off-screen.
        new_blocks = []
        for block in self.blocks:
            block.y += 1
            if block.y < self.height:
                new_blocks.append(block)
        self.blocks = new_blocks
        # Increase score as time passes.
        self.score += 1

    def check_collision(self):
        # Return True if any block occupies the same position as the player.
        for block in self.blocks:
            if block.x == self.player_x and block.y == self.player_y:
                return True
        return False


def draw_border(stdscr, width, height):
    # Draw top and bottom borders.
    for x in range(width):
        try:
            stdscr.addch(0, x, ord("-"))
            stdscr.addch(height - 1, x, ord("-"))
        except curses.error:
            pass
    # Draw left and right borders.
    for y in range(height):
        try:
            stdscr.addch(y, 0, ord("|"))
            stdscr.addch(y, width - 1, ord("|"))
        except curses.error:
            pass
    # Draw corners.
    try:
        stdscr.addch(0, 0, ord("+"))
        stdscr.addch(0, width - 1, ord("+"))
        stdscr.addch(height - 1, 0, ord("+"))
        stdscr.addch(height - 1, width - 1, ord("+"))
    except curses.error:
        pass


def draw_game(stdscr, game, outer_width, outer_height):
    stdscr.erase()  # Erase current content without an immediate physical update.
    draw_border(stdscr, outer_width, outer_height)
    # The playable area is inside the border (offset by 1).
    offset_x, offset_y = 1, 1

    # Draw falling blocks.
    for block in game.blocks:
        x = block.x + offset_x
        y = block.y + offset_y
        if x < outer_width - 1 and y < outer_height - 1:
            try:
                stdscr.addch(y, x, ord("#"))
            except curses.error:
                pass

    # Draw the player with a standout attribute.
    player_x = game.player_x + offset_x
    player_y = game.player_y + offset_y
    if player_x < outer_width - 1 and player_y < outer_height - 1:
        try:
            stdscr.addch(player_y, player_x, ord("@"), curses.A_REVERSE)
        except curses.error:
            pass

    # Draw the score on the top border.
    score_str = f" Score: {game.score} "
    for i, ch in enumerate(score_str):
        if 2 + i < outer_width - 1:
            try:
                stdscr.addch(0, 2 + i, ord(ch))
            except curses.error:
                pass

    stdscr.noutrefresh()  # Mark the window for refresh.
    curses.doupdate()  # Update the physical screen all at once.


def put_centered(stdscr, msg):
    height, width = stdscr.getmaxyx()
    start_x = (width - len(msg)) // 2
    start_y = height // 2
    try:
        stdscr.addstr(start_y, start_x, msg)
    except curses.error:
        pass
    stdscr.noutrefresh()
    curses.doupdate()


def main(stdscr):
    # Initialize colors.
    curses.start_color()
    curses.use_default_colors()
    # Set the background to use the terminal's default.
    stdscr.bkgd(" ", curses.A_NORMAL)

    curses.curs_set(0)  # Hide the cursor.
    stdscr.nodelay(True)  # Non-blocking input.
    stdscr.timeout(0)
    stdscr.erase()

    outer_height, outer_width = stdscr.getmaxyx()
    # Define the playable area inside a 1-character border.
    playable_width = outer_width - 2
    playable_height = outer_height - 2

    game = Game(playable_width, playable_height)
    last_tick = time.time()

    while True:
        # Handle input.
        ch = stdscr.getch()
        if ch != -1:
            if ch == curses.KEY_LEFT:
                if game.player_x > 0:
                    game.player_x -= 1
            elif ch == curses.KEY_RIGHT:
                if game.player_x < game.width - 1:
                    game.player_x += 1
            elif ch == ord("q") or ch == 27:  # 27 is Escape.
                break

        # Update game state on tick.
        now = time.time()
        if now - last_tick >= TICK_RATE:
            game.update()
            if game.check_collision():
                break
            last_tick = now
            draw_game(stdscr, game, outer_width, outer_height)

        time.sleep(0.005)  # Slight delay to reduce CPU usage.

    # Game over: clear the screen and display a final score message.
    stdscr.erase()
    # Reset the background to default before exiting.
    stdscr.bkgd(" ", curses.A_NORMAL)
    msg = f"Game Over! Final Score: {game.score}. Press any key to exit."
    put_centered(stdscr, msg)
    stdscr.nodelay(False)
    stdscr.getch()


if __name__ == "__main__":
    curses.wrapper(main)
