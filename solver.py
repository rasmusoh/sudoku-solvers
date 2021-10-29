import copy

n = 9
c = 3


def run_tests():
    with open('./sudoku.csv') as f:
        next(f)
        for i, line in enumerate(f):
            if i % 1000 == 0:
                print(f'solving at {i}')
            [start, end] = line.split(',')
            try:
                expected = parse_board(end.rstrip('\n'))
                board = solve(parse_board(start))
                if expected != board:
                    print(f'discrepancy at: {i}')
                    print('expected:')
                    print(print_board(expected))
                    print('actual:')
                    print(print_board(board))
                    return
            except:
                print(f'error at: {i}')
                return


def solve(board):
    if not propagate_all(board):
        raise ValueError('board not solvable from starting position')
    elif is_solved(board):
        return board
    hypotheticals = get_next_hypotheticals(board)
    while len(hypotheticals) > 0:
        hyp = hypotheticals.pop()
        if not propagate_all(hyp):
            continue
        elif is_solved(hyp):
            return hyp
        else:
            hypotheticals.extend(get_next_hypotheticals(hyp))
    raise ValueError('Ran out of hypotheticals?')


def get_next_hypotheticals(board):
    minAlts = 10
    minX, minY = -1, -1
    for y, line in enumerate(board):
        for x, cell in enumerate(line):
            alts = len(cell)
            if alts > 1 and alts < minAlts:
                minAlts = alts
                minX, minY = x, y
    hyps = []
    for alt in board[minY][minX]:
        hyp = copy.deepcopy(board)
        hyp[minY][minX] = {alt}
        hyps.append(hyp)
    return hyps


def is_solved(board):
    return all([all([len(cell) == 1 for cell in line]) for line in board])


def propagate_all(board):
    determined = {(i//n, i % n)
                  for i, cell in enumerate([item for line in board for item in line]) if len(cell) == 1}
    while len(determined) > 0:
        y, x = determined.pop()
        if not propagate(y, x, board, determined):
            return False
    return True


def propagate(determinedY, determinedX, board, determined):
    if (len(board[determinedY][determinedX]) != 1):
        raise ValueError(
            f'{determinedY}, {determinedX} not determined, cannot propagate')
    value = next(iter(board[determinedY][determinedX]))
    for y, x in get_neighbors(determinedY, determinedX):
        already_determined = len(board[y][x]) == 1
        board[y][x].discard(value)
        if len(board[y][x]) == 0:
            return False
        elif not already_determined and len(board[y][x]) == 1:
            determined.add((y, x))
    return True


def get_neighbors(y, x):
    horizontal = {(y, i) for i in range(n)}
    vertical = {(i, x) for i in range(n)}
    cell = {(c*(y//c)+i, c*(x//c)+j) for i in range(c) for j in range(c)}
    return (horizontal | vertical | cell) - {(y, x)}


def parse_board(board):
    parsed = [parse_line(board[i:i+n]) for i in range(0, len(board), n)]
    return parsed


def parse_line(line):
    return [set_all() if char == '0' else {int(char)} for char in line]


def set_all():
    return {1, 2, 3, 4, 5, 6, 7, 8, 9}


def print_board(board):
    print('\n')
    for line in board:
        print(' | '.join(
            [str(next(iter(s))) if len(s) == 1 else '-' for s in line]))
    print('\n')


run_tests()
