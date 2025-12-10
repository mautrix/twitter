package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000F\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB'\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\u000e\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0002¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0018\u0010\u000e\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0018\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u000fJ4\u0010\u0011\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\u0010\b\u0002\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0014\u001a\u00020\u0013HÖ\u0001¢\u0006\u0004\b\u0014\u0010\u0015J\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001c\u001a\u00020\u001b2\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001c\u0010\u001dR\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001eR\u001c\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTree;", "Lcom/bendb/thrifty/a;", "", "Lcom/x/dmv2/thriftjava/RatchetTreeLeaf;", "leaves", "Lcom/x/dmv2/thriftjava/RatchetTreeParent;", "parents", "<init>", "(Ljava/util/List;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "copy", "(Ljava/util/List;Ljava/util/List;)Lcom/x/dmv2/thriftjava/RatchetTree;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Companion", "RatchetTreeAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class RatchetTree implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List leaves;

    @JvmField
    @InterfaceC88465b
    public final List parents;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new RatchetTreeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RatchetTree$RatchetTreeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/RatchetTree;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/RatchetTree;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/RatchetTree;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class RatchetTreeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public RatchetTree m85582read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            ArrayList arrayList2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new RatchetTree(arrayList, arrayList2);
                }
                int i = 0;
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 15) {
                        int i2 = protocol.mo14130a2().f38395b;
                        ArrayList arrayList3 = new ArrayList(i2);
                        while (i < i2) {
                            arrayList3.add((RatchetTreeParent) RatchetTreeParent.ADAPTER.read(protocol));
                            i++;
                        }
                        arrayList2 = arrayList3;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 15) {
                    int i3 = protocol.mo14130a2().f38395b;
                    ArrayList arrayList4 = new ArrayList(i3);
                    while (i < i3) {
                        arrayList4.add((RatchetTreeLeaf) RatchetTreeLeaf.ADAPTER.read(protocol));
                        i++;
                    }
                    arrayList = arrayList4;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a RatchetTree struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("RatchetTree");
            if (struct.leaves != null) {
                protocol.mo14136v3("leaves", 1, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.leaves.size());
                Iterator it = struct.leaves.iterator();
                while (it.hasNext()) {
                    RatchetTreeLeaf.ADAPTER.write(protocol, (RatchetTreeLeaf) it.next());
                }
            }
            if (struct.parents != null) {
                protocol.mo14136v3("parents", 2, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.parents.size());
                Iterator it2 = struct.parents.iterator();
                while (it2.hasNext()) {
                    RatchetTreeParent.ADAPTER.write(protocol, (RatchetTreeParent) it2.next());
                }
            }
            protocol.mo14134i0();
        }
    }

    public RatchetTree(@InterfaceC88465b List list, @InterfaceC88465b List list2) {
        this.leaves = list;
        this.parents = list2;
    }

    public static /* synthetic */ RatchetTree copy$default(RatchetTree ratchetTree, List list, List list2, int i, Object obj) {
        if ((i & 1) != 0) {
            list = ratchetTree.leaves;
        }
        if ((i & 2) != 0) {
            list2 = ratchetTree.parents;
        }
        return ratchetTree.copy(list, list2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getLeaves() {
        return this.leaves;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getParents() {
        return this.parents;
    }

    @InterfaceC88464a
    public final RatchetTree copy(@InterfaceC88465b List leaves, @InterfaceC88465b List parents) {
        return new RatchetTree(leaves, parents);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof RatchetTree)) {
            return false;
        }
        RatchetTree ratchetTree = (RatchetTree) other;
        return Intrinsics.m65267c(this.leaves, ratchetTree.leaves) && Intrinsics.m65267c(this.parents, ratchetTree.parents);
    }

    public int hashCode() {
        List list = this.leaves;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        List list2 = this.parents;
        return iHashCode + (list2 != null ? list2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "RatchetTree(leaves=" + this.leaves + ", parents=" + this.parents + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}