# Legend
$M \dots \text{mass of body}$
$\vec{F} \dots \text{force acting on body}$
$\vec{I} \dots \text{inertia of body}$
$\vec{P} \dots \text{position of body}$
$t' \dots \text{delta time}$
# Approaches
As the differential equations for 3 or more bodies are not solvable, the simulation approximates the behaviour by simulating in discrete timesteps ($t'$).
## Primitive

$\vec{I'} = \vec{I} + \frac{\vec{F}}{M}t'$
$\vec{P'} = \vec{P} + \vec{I'} t'$

Calculates forces and changes momentum by multiplying the force-mass quotient by delta time. The new position is calculates by applying the new momentum for the duration $t'$.

### Effects
This simple approach results in semmingly correct behaviour. As celestial body movement is chaotic the behaviour might seem correct over a short timespan but results in wrong movement in the long run. No simulation method can avoid this but the error can be made smaller.

The most noticable effect this approach has: 'drifting'. The distance between two objects orbiting aroud each other (e.g. the earth and moon pair) will slowly increase.

![[drift.jpg]]

To get rid of this "drift" / reduce it, a different approach is needed.

## First derivative
$F_n(t) = \sum_{k\ne n} G\frac{m_n m_k}{d_{n,k}(t)^2}$
$I_n(t) = \int_{0}^{t}F_n(\tau)d\tau$
$P_n(t) = \frac{1}{m_n}\int_{0}^{t}I_n(\tau)d\tau$ 

In this approach we want to make the Position $P_n$ dependent on the change of $F_n$ which is only dependent in the change of the distance $d_{n,k}$.

$F_n'(t) = f(t) = \sum_{k \ne n} G m_n m_k \cdot -\frac{1}{d_{n,k}(t)^4} \cdot 2d_{n,k}(t) \cdot d_{n,k}'(t) = \sum_{k \ne n} G m_n m_k \cdot - \frac{2d_{n,k}'(t)}{d_{n,k}(t)^3}$
$I_n = \int_0^t F_n(\tau) d\tau = \int_0^t \int_0^\tau f(\tau') d\tau' d\tau$

$\vec{P'} = \vec{P} + \vec{I} t' + \frac{\vec{F}}{2M} t'^2$
$\vec{I'} = \vec{I} + \frac{\vec{F}}{M} t'$

This approach updates the position by adding a linear movement for duration $t'$ and a parabolic movement dependent on the force. In direct comparison to the primitive approach it only assumes the force to be constant not the inertia during $t'$.

The Inertia is updated the same way as in the primitive approach.

### Effects
Not testet yet.
Expected:
- drift (hopefully less than with primitive approach)
