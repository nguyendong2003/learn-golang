// Header component — sticky navbar, transparent on home hero
function Header({ navigate, currentPage }) {
  const [scrolled, setScrolled] = useState(false);
  const [notifOpen, setNotifOpen] = useState(false);
  const [saved, setSaved] = useState(false);
  const isHome = currentPage === 'home';
  const transp = isHome && !scrolled;

  useEffect(() => {
    const fn = () => setScrolled(window.scrollY > 60);
    window.addEventListener('scroll', fn);
    return () => window.removeEventListener('scroll', fn);
  }, []);

  const navLinks = [
    { label:'Tour', page:'search', p:{ type:'tour' } },
    { label:'Khách sạn', page:'search', p:{ type:'hotel' } },
    { label:'Combo', page:'search', p:{ type:'combo' } },
    { label:'Blog', page:'blog', p:{} },
    { label:'Khuyến mãi', page:'deals', p:{} },
  ];

  const notifs = [
    { icon:'✅', text:'Booking BK260501 đã được xác nhận', time:'2 giờ trước', unread:true },
    { icon:'💬', text:'Halong Star Travel vừa trả lời câu hỏi của bạn', time:'1 ngày trước', unread:true },
    { icon:'🎁', text:'Bạn có voucher SUMMER30 — Giảm 30%', time:'3 ngày trước', unread:false },
  ];

  const hdr = {
    position:'fixed', top:0, left:0, right:0, zIndex:1000,
    background: transp ? 'transparent' : 'rgba(255,255,255,0.97)',
    backdropFilter: transp ? 'none' : 'blur(14px)',
    borderBottom: transp ? 'none' : '1px solid var(--border)',
    boxShadow: transp ? 'none' : 'var(--sh-xs)',
    transition:'all .3s ease',
  };

  const logoColor = transp ? '#fff' : 'var(--primary)';
  const accentColor = transp ? 'rgba(255,255,255,0.65)' : 'var(--accent)';
  const linkColor = transp ? 'rgba(255,255,255,0.88)' : 'var(--text2)';
  const linkHov = transp ? '#fff' : 'var(--primary)';

  return (
    <header style={hdr}>
      <div className="container" style={{ height:68, display:'flex', alignItems:'center', gap:0 }}>

        {/* Logo */}
        <button onClick={() => navigate('home')} style={{ display:'flex', alignItems:'center', gap:8, background:'none', border:'none', cursor:'pointer', padding:0, marginRight:36, flexShrink:0 }}>
          <div style={{ width:34, height:34, borderRadius:10, background: transp ? 'rgba(255,255,255,0.18)' : 'var(--primary)', display:'flex', alignItems:'center', justifyContent:'center' }}>
            <Ico.pin />
          </div>
          <span style={{ fontSize:19, fontWeight:800, letterSpacing:-.5, color:logoColor }}>
            Travel<span style={{ color:accentColor }}>Viet</span>
          </span>
        </button>

        {/* Nav */}
        <nav style={{ display:'flex', gap:2, flex:1 }}>
          {navLinks.map(link => (
            <NavLink key={link.label} label={link.label} color={linkColor} hov={linkHov}
              active={currentPage === link.page}
              transp={transp}
              onClick={() => navigate(link.page, link.p)} />
          ))}
        </nav>

        {/* Right */}
        <div style={{ display:'flex', alignItems:'center', gap:6 }}>

          {/* Bell */}
          <div style={{ position:'relative' }}>
            <button onClick={() => setNotifOpen(!notifOpen)} style={{ width:38, height:38, borderRadius:'50%', border:'none', cursor:'pointer', display:'flex', alignItems:'center', justifyContent:'center', background: transp ? 'rgba(255,255,255,0.12)' : 'var(--border-lt)', color: transp ? '#fff' : 'var(--text2)', position:'relative', transition:'all .15s' }}>
              <Ico.bell />
              <span style={{ position:'absolute', top:7, right:7, width:8, height:8, borderRadius:'50%', background:'var(--accent)', border:'2px solid ' + (transp ? 'transparent' : '#fff') }}></span>
            </button>
            {notifOpen && (
              <div className="dropdown-anim" style={{ position:'absolute', top:'calc(100% + 10px)', right:0, width:340, background:'#fff', borderRadius:var_r_lg, boxShadow:'var(--sh-xl)', border:'1px solid var(--border)', zIndex:200, overflow:'hidden' }}>
                <div style={{ padding:'14px 18px 10px', display:'flex', justifyContent:'space-between', alignItems:'center', borderBottom:'1px solid var(--border-lt)' }}>
                  <span style={{ fontWeight:700, fontSize:15 }}>Thông báo</span>
                  <button style={{ fontSize:12, color:'var(--primary)', background:'none', border:'none', cursor:'pointer', fontWeight:600 }}>Đọc tất cả</button>
                </div>
                {notifs.map((n,i) => (
                  <div key={i} style={{ padding:'12px 18px', display:'flex', gap:12, alignItems:'flex-start', background: n.unread ? 'var(--primary-xlt)' : 'transparent', borderBottom:'1px solid var(--border-lt)', cursor:'pointer' }}>
                    <span style={{ fontSize:20, flexShrink:0 }}>{n.icon}</span>
                    <div style={{ flex:1 }}>
                      <div style={{ fontSize:13, fontWeight: n.unread ? 600 : 400, lineHeight:1.4 }}>{n.text}</div>
                      <div style={{ fontSize:11, color:'var(--text3)', marginTop:3 }}>{n.time}</div>
                    </div>
                    {n.unread && <div style={{ width:7, height:7, borderRadius:'50%', background:'var(--primary)', flexShrink:0, marginTop:5 }}></div>}
                  </div>
                ))}
                <div style={{ padding:'10px 18px', textAlign:'center' }}>
                  <button style={{ fontSize:13, color:'var(--primary)', background:'none', border:'none', cursor:'pointer', fontWeight:600 }}>Xem tất cả thông báo</button>
                </div>
              </div>
            )}
          </div>

          {/* Language */}
          <button style={{ height:34, padding:'0 12px', borderRadius:'var(--r-full)', border: transp ? '1.5px solid rgba(255,255,255,0.35)' : '1.5px solid var(--border)', background:'transparent', color: transp ? '#fff' : 'var(--text2)', fontSize:13, fontWeight:600, cursor:'pointer', display:'flex', alignItems:'center', gap:5, transition:'all .15s' }}>
            <Ico.globe /> VI
          </button>

          {/* Đăng nhập / Account */}
          <button onClick={() => navigate('account')} style={{ height:38, padding:'0 16px', borderRadius:'var(--r-full)', border:'none', background: transp ? 'rgba(255,255,255,0.18)' : 'var(--primary-lt)', color: transp ? '#fff' : 'var(--primary)', fontSize:13, fontWeight:700, cursor:'pointer', display:'flex', alignItems:'center', gap:8, transition:'all .15s' }}>
            <div style={{ width:26, height:26, borderRadius:'50%', background: transp ? 'rgba(255,255,255,0.3)' : 'var(--primary)', color:'#fff', display:'flex', alignItems:'center', justifyContent:'center', fontSize:11, fontWeight:800 }}>NT</div>
            Nguyễn Tuấn
          </button>

        </div>
      </div>
      {notifOpen && <div onClick={() => setNotifOpen(false)} style={{ position:'fixed', inset:0, zIndex:199 }} />}
    </header>
  );
}

const var_r_lg = 'var(--r-lg)';

function NavLink({ label, color, hov, active, transp, onClick }) {
  const [hover, setHover] = useState(false);
  return (
    <button onClick={onClick}
      onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}
      style={{ background:'none', border:'none', cursor:'pointer', padding:'6px 12px', borderRadius:8, fontWeight: active ? 700 : 500, fontSize:14, color: hover || active ? hov : color, transition:'all .15s', position:'relative' }}>
      {label}
      {active && !transp && <div style={{ position:'absolute', bottom:2, left:'50%', transform:'translateX(-50%)', width:20, height:2.5, borderRadius:2, background:'var(--primary)' }} />}
    </button>
  );
}

Object.assign(window, { Header });
