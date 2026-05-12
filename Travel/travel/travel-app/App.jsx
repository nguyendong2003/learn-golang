// App.jsx — Main router + ChatBot bubble
function ChatBot() {
  const [open, setOpen] = useState(false);
  const [msgs, setMsgs] = useState([
    { from:'bot', text:'Xin chào! Tôi là trợ lý AI của TravelViet. Tôi có thể giúp bạn tìm tour, khách sạn hoặc tư vấn lịch trình du lịch. Bạn muốn đi đâu?' }
  ]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const endRef = useRef(null);

  useEffect(() => {
    if (open && endRef.current) endRef.current.scrollTop = endRef.current.scrollHeight;
  }, [msgs, open]);

  async function send() {
    if (!input.trim() || loading) return;
    const userMsg = input.trim();
    setInput('');
    setMsgs(m => [...m, { from:'user', text: userMsg }]);
    setLoading(true);
    try {
      const reply = await window.claude.complete({
        messages: [
          { role:'user', content: 'Bạn là trợ lý AI của TravelViet — nền tảng đặt tour & khách sạn nội địa Việt Nam. Hãy trả lời ngắn gọn (dưới 80 từ), thân thiện, bằng tiếng Việt. Câu hỏi: ' + userMsg }
        ]
      });
      setMsgs(m => [...m, { from:'bot', text: reply }]);
    } catch {
      setMsgs(m => [...m, { from:'bot', text: 'Xin lỗi, tôi đang gặp sự cố. Vui lòng thử lại hoặc gọi hotline 1800 6006!' }]);
    }
    setLoading(false);
  }

  const quickReplies = ['Gợi ý tour Phú Quốc','KS Đà Nẵng giá tốt','Combo tiết kiệm nhất'];

  return (
    <>
      {open && (
        <div style={{ position:'fixed', bottom:90, right:24, width:340, height:480, background:'#fff', borderRadius:'var(--r-xl)', boxShadow:'var(--sh-xl)', zIndex:9999, display:'flex', flexDirection:'column', overflow:'hidden', border:'1.5px solid var(--border)' }}>
          <div style={{ background:'var(--primary)', padding:'16px 18px', display:'flex', alignItems:'center', gap:12 }}>
            <div style={{ width:36, height:36, borderRadius:'50%', background:'rgba(255,255,255,.2)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:18 }}>🤖</div>
            <div style={{ flex:1 }}>
              <div style={{ fontWeight:700, color:'#fff', fontSize:14 }}>TravelBot AI</div>
              <div style={{ fontSize:12, color:'rgba(255,255,255,.7)', display:'flex', alignItems:'center', gap:4 }}>
                <span style={{ width:6, height:6, borderRadius:'50%', background:'#4ade80', display:'inline-block' }}></span> Đang hoạt động
              </div>
            </div>
            <button onClick={() => setOpen(false)} style={{ background:'rgba(255,255,255,.15)', border:'none', color:'#fff', width:28, height:28, borderRadius:'50%', cursor:'pointer', display:'flex', alignItems:'center', justifyContent:'center', fontSize:14 }}>✕</button>
          </div>
          <div ref={endRef} style={{ flex:1, overflow:'auto', padding:'14px 16px', display:'flex', flexDirection:'column', gap:10, background:'var(--bg)' }}>
            {msgs.map((m, i) => (
              <div key={i} style={{ display:'flex', justifyContent: m.from==='user' ? 'flex-end' : 'flex-start' }}>
                {m.from === 'bot' && <div style={{ width:28, height:28, borderRadius:'50%', background:'var(--primary)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:13, marginRight:8, flexShrink:0, marginTop:2 }}>🤖</div>}
                <div style={{ maxWidth:'78%', padding:'10px 14px', borderRadius: m.from==='user' ? '18px 18px 4px 18px' : '18px 18px 18px 4px', background: m.from==='user' ? 'var(--primary)' : '#fff', color: m.from==='user' ? '#fff' : 'var(--text)', fontSize:13, lineHeight:1.6, boxShadow:'var(--sh-xs)' }}>{m.text}</div>
              </div>
            ))}
            {loading && (
              <div style={{ display:'flex', alignItems:'center', gap:8 }}>
                <div style={{ width:28, height:28, borderRadius:'50%', background:'var(--primary)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:13 }}>🤖</div>
                <div style={{ background:'#fff', borderRadius:'18px 18px 18px 4px', padding:'10px 14px', display:'flex', gap:5, alignItems:'center' }}>
                  {[0,1,2].map(i => <div key={i} style={{ width:7, height:7, borderRadius:'50%', background:'var(--text3)', animation:'pulse 1.2s ease infinite', animationDelay: i*0.2+'s' }} />)}
                </div>
              </div>
            )}
          </div>
          {msgs.length <= 2 && (
            <div style={{ padding:'8px 12px', display:'flex', gap:6, flexWrap:'wrap', borderTop:'1px solid var(--border-lt)', background:'#fff' }}>
              {quickReplies.map(q => (
                <button key={q} onClick={() => setInput(q)} style={{ padding:'5px 12px', borderRadius:'var(--r-full)', border:'1.5px solid var(--primary)', background:'var(--primary-xlt)', color:'var(--primary)', fontSize:12, fontWeight:600, cursor:'pointer' }}>{q}</button>
              ))}
            </div>
          )}
          <div style={{ padding:'12px 14px', borderTop:'1px solid var(--border-lt)', display:'flex', gap:8, background:'#fff' }}>
            <input value={input} onChange={e => setInput(e.target.value)} onKeyDown={e => e.key==='Enter' && send()}
              placeholder="Nhập câu hỏi..." style={{ flex:1, border:'1.5px solid var(--border)', borderRadius:'var(--r-full)', padding:'8px 14px', fontSize:13, outline:'none' }}
              onFocus={e => e.target.style.borderColor='var(--primary)'} onBlur={e => e.target.style.borderColor='var(--border)'} />
            <button onClick={send} disabled={!input.trim() || loading}
              style={{ width:36, height:36, borderRadius:'50%', background:'var(--primary)', border:'none', color:'#fff', cursor:'pointer', display:'flex', alignItems:'center', justifyContent:'center', opacity: (!input.trim()||loading) ? 0.5 : 1 }}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/></svg>
            </button>
          </div>
        </div>
      )}
      <button onClick={() => setOpen(!open)}
        style={{ position:'fixed', bottom:24, right:24, width:58, height:58, borderRadius:'50%', background: open ? 'var(--text)' : 'var(--primary)', border:'none', color:'#fff', cursor:'pointer', boxShadow:'var(--sh-lg)', zIndex:9999, display:'flex', alignItems:'center', justifyContent:'center', transition:'all .25s', fontSize: open ? 20 : 22 }}>
        {open ? '✕' : <Ico.chat />}
        {!open && <span style={{ position:'absolute', top:10, right:10, width:10, height:10, borderRadius:'50%', background:'var(--accent)', border:'2px solid #fff' }} />}
      </button>
    </>
  );
}

function App() {
  const [page, setPage] = useState('home');
  const [params, setParams] = useState({});

  function navigate(pageName, p = {}) {
    setPage(pageName);
    setParams(p);
    window.scrollTo({ top:0, behavior:'smooth' });
  }

  return (
    <div>
      <Header navigate={navigate} currentPage={page} />

      {page === 'home'         && <HomePage navigate={navigate} />}
      {page === 'search'       && <SearchPage navigate={navigate} params={params} />}
      {page === 'tour-detail'  && <TourDetailPage navigate={navigate} params={params} />}
      {page === 'hotel-detail' && <HotelDetailPage navigate={navigate} params={params} />}
      {page === 'combo-detail' && <ComboDetailPage navigate={navigate} params={params} />}
      {page === 'booking'      && <BookingPage navigate={navigate} params={params} />}
      {page === 'account'      && <AccountPage navigate={navigate} />}
      {page === 'deals'        && <DealsPage navigate={navigate} />}

      {['blog'].includes(page) && (
        <div style={{ paddingTop:68, minHeight:'50vh', display:'flex', alignItems:'center', justifyContent:'center', flexDirection:'column', gap:12, color:'var(--text3)' }}>
          <div style={{ fontSize:44 }}>🚧</div>
          <div style={{ fontWeight:700, fontSize:18 }}>Trang đang phát triển</div>
          <button onClick={() => navigate('home')} className="btn btn-primary btn-sm" style={{ marginTop:8 }}>← Về trang chủ</button>
        </div>
      )}

      <ChatBot />
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(<App />);
