package rush

type Solver struct {
	board  *Board
	target int
	memo   *Memo
	path   []Move
	moves  [][]Move
}

type Solution struct {
	Solvable bool
	Moves    []Move
	NumMoves int
	NumSteps int
	Depth    int
	MemoSize int
	MemoHits uint64
}

func NewSolver(board *Board) *Solver {
	solver := Solver{}
	solver.board = board
	solver.target = board.Target()
	solver.memo = NewMemo()
	return &solver
}

func (solver *Solver) isSolved() bool {
	return solver.board.Pieces[0].Position == solver.target
}

func (solver *Solver) search(depth, maxDepth int) bool {
	height := maxDepth - depth
	if height == 0 {
		return solver.isSolved()
	}

	board := solver.board
	if !solver.memo.Add(board.MemoKey(), height) {
		return false
	}

	// count occupied squares between primary piece and target
	primary := board.Pieces[0]
	i0 := primary.Position + primary.Size
	i1 := solver.target + primary.Size - 1
	minMoves := 0
	for i := i0; i <= i1; i++ {
		if board.occupied[i] {
			minMoves++
		}
	}
	if minMoves >= height {
		return false
	}

	buf := &solver.moves[depth]
	*buf = board.Moves(*buf)
	for _, move := range *buf {
		board.DoMove(move)
		solved := solver.search(depth+1, maxDepth)
		board.UndoMove(move)
		if solved {
			solver.path[depth] = move
			return true
		}
	}
	return false
}

func (solver *Solver) Solve() Solution {
	board := solver.board
	memo := solver.memo

	if err := board.Validate(); err != nil {
		return Solution{}
	}

	if solver.isSolved() {
		return Solution{Solvable: true}
	}

	if board.Impossible() {
		return Solution{}
	}

	previousMemoSize := 0
	cutoff := board.Width - board.Pieces[0].Size
	for i := 1; ; i++ {
		solver.path = make([]Move, i)
		solver.moves = make([][]Move, i)
		if solver.search(0, i) {
			moves := solver.path
			steps := 0
			for _, move := range moves {
				steps += move.AbsSteps()
			}
			return Solution{
				Solvable: true,
				Moves:    moves,
				NumMoves: len(moves),
				NumSteps: steps,
				Depth:    i,
				MemoSize: memo.Size(),
				MemoHits: memo.Hits(),
			}
		}
		memoSize := memo.Size()
		if i > cutoff && memoSize == previousMemoSize {
			return Solution{
				Depth:    i,
				MemoSize: memo.Size(),
				MemoHits: memo.Hits(),
			}
		}
		previousMemoSize = memoSize
	}
}