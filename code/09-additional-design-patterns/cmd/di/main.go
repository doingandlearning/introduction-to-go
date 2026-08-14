// Command di demonstrates manual dependency injection: a larger
// dependency graph than a single constructor argument. OrderHandler
// depends on OrderService, which depends on two interfaces - Repository
// and Notifier - all wired together by hand in main, with no framework,
// no reflection, no container.
package main

import "fmt"

// Repository and Notifier are the two dependencies OrderService needs.
type Repository interface {
	Save(order string) error
}

type Notifier interface {
	Notify(message string)
}

type OrderService struct {
	repo     Repository
	notifier Notifier
}

func NewOrderService(repo Repository, notifier Notifier) *OrderService {
	return &OrderService{repo: repo, notifier: notifier}
}

func (s *OrderService) PlaceOrder(order string) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}
	s.notifier.Notify("order placed: " + order)
	return nil
}

// OrderHandler previews the Handler layer from Topic 8's Handler ->
// Service -> Repository architecture. It depends only on OrderService,
// never on Repository or Notifier directly.
type OrderHandler struct {
	svc *OrderService
}

func NewOrderHandler(svc *OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) HandlePlaceOrder(order string) {
	if err := h.svc.PlaceOrder(order); err != nil {
		fmt.Println("failed:", err)
		return
	}
	fmt.Println("handled:", order)
}

// --- real implementations ---

type InMemoryRepository struct {
	saved []string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) Save(order string) error {
	r.saved = append(r.saved, order)
	return nil
}

type ConsoleNotifier struct{}

func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (ConsoleNotifier) Notify(message string) {
	fmt.Println("[notify]", message)
}

// --- fake/test-double implementations ---
// No interface changed, no mocking framework involved - these are just
// ordinary types that satisfy Repository and Notifier.

type FakeRepository struct {
	saved []string
}

func (r *FakeRepository) Save(order string) error {
	r.saved = append(r.saved, order)
	return nil
}

type FakeNotifier struct {
	messages []string
}

func (n *FakeNotifier) Notify(message string) {
	n.messages = append(n.messages, message)
}

func main() {
	// Manual DI: every wire is a plain constructor call, visible top to
	// bottom, right here in main. Compare this to Spring's @Autowired -
	// there, the graph is resolved implicitly at runtime by a container.
	// Here, you are the container, and it's a design choice, not a gap.
	fmt.Println("-- real wiring --")
	repo := NewInMemoryRepository()
	notifier := NewConsoleNotifier()
	svc := NewOrderService(repo, notifier)
	handler := NewOrderHandler(svc)
	handler.HandlePlaceOrder("2x conference badge")

	// Same OrderService and OrderHandler code, wired a second time with
	// fakes. Nothing about either type changed to make this possible -
	// they only ever depended on interfaces, never on a concrete type.
	fmt.Println("-- fake wiring (what a test would do) --")
	fakeRepo := &FakeRepository{}
	fakeNotifier := &FakeNotifier{}
	testSvc := NewOrderService(fakeRepo, fakeNotifier)
	testHandler := NewOrderHandler(testSvc)
	testHandler.HandlePlaceOrder("1x workshop ticket")

	fmt.Println("fake repo saved:", fakeRepo.saved)
	fmt.Println("fake notifier messages:", fakeNotifier.messages)
}
