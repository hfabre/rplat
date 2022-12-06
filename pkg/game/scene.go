package game

type Scene interface {
	Init()
	UpdateInputs()
	HandleEvents()
	ClearInputs()
	Update(float32)
	Draw(float64)
	End()
	ShouldExit() bool
}

type SceneManager struct {
	scenes           map[string]Scene
	currentSceneName string
}

func (sm *SceneManager) SwapScene(scene string) {
	oldScene := sm.CurrentScene()
	sm.currentSceneName = scene
	oldScene.End()
	sm.CurrentScene().Init()
}

func (sm SceneManager) UpdateInputs() {
	sm.CurrentScene().UpdateInputs()
}

func (sm SceneManager) ClearInputs() {
	sm.CurrentScene().ClearInputs()
}

func (sm SceneManager) HandleEvents() {
	sm.CurrentScene().HandleEvents()
}

func (sm SceneManager) Update(dt float32) {
	sm.CurrentScene().Update(dt)
}

func (sm SceneManager) Draw(factor float64) {
	sm.CurrentScene().Draw(factor)
}

func (sm SceneManager) CurrentScene() Scene {
	return sm.scenes[sm.currentSceneName]
}

func (sm SceneManager) ShouldExit() bool {
	return sm.CurrentScene().ShouldExit()
}
